/**
 * A OpenTraced server for a go service that implements the commandLine interface.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	pb "../command_line"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"../../../../go/otgrpc"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
)

const (
	port = ":50051"
)

var accessToken = flag.String("access_token", "", "your LightStep access token")

// server is used to implement helloworld.CommandLineServer.
type commandLineServer struct {
	Inventory map[string]int
}

func (s *commandLineServer) Echo(ctx context.Context, in *pb.CommandRequest) (*pb.CommandResponse, error) {
	return &pb.CommandResponse{in.Text}, nil
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
	tracerOpts.Tags[lightstep.ComponentNameKey] = "go.trivial-server"
	tracer := lightstep.NewTracer(tracerOpts)

	// Set up a connection to the server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())))
	pb.RegisterCommandLineServer(s, &commandLineServer{})
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
