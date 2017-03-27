/**
 * A OpenTraced client for a go service that implements the commandLine interface.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pb "../command_line"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"../../../../go/otgrpc"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
)

const (
	address = "localhost:50051"
)

var accessToken = flag.String("access_token", "", "your LightStep access token")

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
	tracerOpts.Tags[lightstep.ComponentNameKey] = "go.trivial-client"
	tracer := lightstep.NewTracer(tracerOpts)

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCommandLineClient(conn)

	// ctx, cancel := context.WithCancel(context.Background())
	// cancel()
	// resp, err := c.Echo(ctx, &pb.CommandRequest{"hello, hello"})
	resp, err := c.Echo(context.Background(), &pb.CommandRequest{"hello, hello"})
	if err != nil {
		// log.Fatalf("Echo failure: %v", err)
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(resp.Text)
	}

	// Force a flush of the tracer.
	err = lightstep.FlushLightStepTracer(tracer)
	if err != nil {
		panic(err)
	}
}
