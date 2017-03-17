/**
 * A OpenTraced client for a go service that implements the store interface.
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "../store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"../../../../go/otgrpc"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
)

const (
	address = "localhost:50051"
	help    = `Enter commands to interact with the store service:

    stock_item     Stock a single item.
    stock_items    Stock one or more items.
    sell_item      Sell a single item.
    sell_items     Sell one or more items.
    inventory      List the store's inventory.
    query_item     Query the inventory for a single item.
    query_items    Query the inventory for one or more items.

Example:
    > stock_item apple
    > stock_items apple milk
    > inventory
    apple   2
    milk    1
`
)

var accessToken = flag.String("access_token", "", "your LightStep access token")

func ExecuteCommand(c pb.StoreClient, cmd string, args []string) error {
	switch cmd {
	case "stock_item":
		if len(args) != 1 {
			fmt.Println("must input a single item")
			return nil
		}
		_, err := c.AddItem(context.Background(), &pb.AddItemRequest{args[0]})
		return err
	case "stock_items":
		if len(args) == 0 {
			fmt.Println("must input at least a single item")
			return nil
		}
		stream, err := c.AddItems(context.Background())
		if err != nil {
			return err
		}
		for _, item := range args {
			if err = stream.Send(&pb.AddItemRequest{item}); err != nil {
				return err
			}
		}
		_, err = stream.CloseAndRecv()
		return err
	case "stock_items_error":
		if len(args) == 0 {
			fmt.Println("must input at least a single item")
			return nil
		}
		stream, err := c.AddItemsError(context.Background())
		if err != nil {
			return err
		}
		for _, item := range args {
			if err = stream.Send(&pb.AddItemRequest{item}); err != nil {
				// Ignore the error
				return nil
			}
		}
		_, err = stream.CloseAndRecv()
		return nil
	case "sell_item":
		if len(args) != 1 {
			fmt.Println("must input a single item")
			return nil
		}
		was_sold, err := c.RemoveItem(context.Background(), &pb.RemoveItemRequest{args[0]})
		if err != nil {
			return err
		}
		if !was_sold.WasSuccessful {
			fmt.Println("unable to sell")
		}
	case "sell_items":
		if len(args) == 0 {
			fmt.Println("must input at least a single item")
			return nil
		}
		stream, err := c.RemoveItems(context.Background())
		if err != nil {
			return err
		}
		for _, item := range args {
			if err = stream.Send(&pb.RemoveItemRequest{item}); err != nil {
				return err
			}
		}
		was_sold, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}
		if !was_sold.WasSuccessful {
			fmt.Println("unable to sell")
		}
	case "inventory":
		if len(args) != 0 {
			fmt.Println("inventory does not take any args")
			return nil
		}
		stream, err := c.ListInventory(context.Background(), &pb.Empty{})
		if err != nil {
			return err
		}
		for {
			item_quantity, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Printf("%s\t%d\n", item_quantity.Name, item_quantity.Count)
		}
	case "inventory_cancel":
		if len(args) != 0 {
			fmt.Println("inventory does not take any args")
			return nil
		}
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := c.ListInventory(ctx, &pb.Empty{})
		if err != nil {
			return err
		}
		for {
			item_quantity, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				// Ignore the error
				return nil
			}
			fmt.Printf("%s\t%d\n", item_quantity.Name, item_quantity.Count)
			cancel()
		}
	case "query_item":
		if len(args) != 1 {
			fmt.Println("must input a single item")
			return nil
		}
		item_quantity, err := c.QueryQuantity(context.Background(), &pb.QueryItemRequest{args[0]})
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%d\n", item_quantity.Name, item_quantity.Count)
	case "query_items":
		if len(args) == 0 {
			fmt.Println("must input at least a single item")
			return nil
		}
		stream, err := c.QueryQuantities(context.Background())
		if err != nil {
			return err
		}
		waitc := make(chan error)
		go func() {
			for {
				item_quantity, err := stream.Recv()
				if err == io.EOF {
					waitc <- nil
					return
				}
				if err != nil {
					waitc <- err
					return
				}
				fmt.Printf("%s\t%d\n", item_quantity.Name, item_quantity.Count)
			}
		}()
		for _, item := range args {
			if err := stream.Send(&pb.QueryItemRequest{item}); err != nil {
				return err
			}
		}
		stream.CloseSend()
		return <-waitc
	default:
		fmt.Printf("unknown command \"%s\"\n", cmd)
	}
	return nil
}

func ReadAndExecute(c pb.StoreClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(help)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("failed to read string: %v", err)
		}
		components := strings.Fields(text)
		if len(components) == 0 {
			continue
		}
		cmd := components[0]
		args := components[1:]
		if err = ExecuteCommand(c, cmd, args); err != nil {
			log.Fatalf("failed to execute command: %v", err)
		}
	}
}

func main() {
	flag.Parse()
	if len(*accessToken) == 0 {
		fmt.Println("You must specify --access_token")
		os.Exit(1)
	}

	// Set up the LightStep tracer
	tracerOpts := lightstep.Options{
		AccessToken: *accessToken,
		UseGRPC:     true,
	}
	tracerOpts.Tags = make(opentracing.Tags)
	tracerOpts.Tags[lightstep.ComponentNameKey] = "go.store-client"
	tracer := lightstep.NewTracer(tracerOpts)

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStoreClient(conn)

	// Read and execute commands for the store.
	ReadAndExecute(c)

	// Force a flush of the tracer.
	err = lightstep.FlushLightStepTracer(tracer)
	if err != nil {
		panic(err)
	}
}
