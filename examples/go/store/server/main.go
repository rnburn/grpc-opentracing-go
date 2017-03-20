/**
 * A OpenTraced server for a go service that implements the store interface.
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "../store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"../../../../go/otgrpc"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

const (
	port = ":50051"
)

var accessToken = flag.String("access_token", "", "your LightStep access token")

// server is used to implement helloworld.StoreServer.
type storeServer struct {
	Inventory map[string]int
}

func (s *storeServer) AddItem(ctx context.Context, in *pb.AddItemRequest) (*pb.Empty, error) {
	s.Inventory[in.Name]++
	return &pb.Empty{}, nil
}

func (s *storeServer) AddItems(stream pb.Store_AddItemsServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Empty{})
		}
		if err != nil {
			return err
		}
		s.Inventory[item.Name]++
	}
}

func (s *storeServer) AddItemsError(stream pb.Store_AddItemsErrorServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Empty{})
		}
		if err != nil {
			return err
		}
		s.Inventory[item.Name]++
		return errors.New("Failed to add item.")
	}
}

func (s *storeServer) RemoveItem(ctx context.Context, in *pb.RemoveItemRequest) (*pb.RemoveItemResponse, error) {
	if quantity := s.Inventory[in.Name]; quantity > 0 {
		s.Inventory[in.Name]--
		return &pb.RemoveItemResponse{true}, nil
	} else {
		return &pb.RemoveItemResponse{false}, nil
	}
}

func (s *storeServer) RemoveItems(stream pb.Store_RemoveItemsServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.RemoveItemResponse{true})
		}
		if err != nil {
			return err
		}
		if quantity := s.Inventory[item.Name]; quantity > 0 {
			s.Inventory[item.Name]--
		} else {
			return stream.SendAndClose(&pb.RemoveItemResponse{false})
		}
	}
}

func (s *storeServer) ListInventory(_ *pb.Empty, stream pb.Store_ListInventoryServer) error {
	for name, quantity := range s.Inventory {
		if err := stream.Send(&pb.QuantityResponse{name, int32(quantity)}); err != nil {
			return err
		}
	}
	return nil
}

func (s *storeServer) QueryQuantity(ctx context.Context, in *pb.QueryItemRequest) (*pb.QuantityResponse, error) {
	quantity := s.Inventory[in.Name]
	return &pb.QuantityResponse{in.Name, int32(quantity)}, nil
}

func (s *storeServer) QueryQuantities(stream pb.Store_QueryQuantitiesServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		quantity := s.Inventory[item.Name]
		if err = stream.Send(&pb.QuantityResponse{item.Name, int32(quantity)}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	if len(*accessToken) == 0 {
		fmt.Println("You must specify --access_token")
		os.Exit(1)
	}

	tracerOpts := lightstep.Options{
		AccessToken: *accessToken,
		UseGRPC:     true,
	}
	tracerOpts.Tags = make(opentracing.Tags)
	tracerOpts.Tags[lightstep.ComponentNameKey] = "go.store-server"
	tracer := lightstep.NewTracer(tracerOpts)

	// Set up a span decorator
	decorator := func(
		span opentracing.Span,
		method string,
		req, resp interface{},
		err error) {
		span.LogFields(
			otlog.String("event", "decoration"),
			otlog.String("method", method),
			otlog.Error(err))
	}

	// Set up a connection to the server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads(), otgrpc.SpanDecorator(decorator))),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer, otgrpc.LogPayloads(), otgrpc.SpanDecorator(decorator))))
	pb.RegisterStoreServer(s, &storeServer{make(map[string]int)})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Force a flush of the tracer
	err = lightstep.FlushLightStepTracer(tracer)
	if err != nil {
		panic(err)
	}
}
