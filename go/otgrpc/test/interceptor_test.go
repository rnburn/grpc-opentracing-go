package interceptor_test

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	// testpb "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc/test/otgrpc_testing"
	"../../otgrpc"
	testpb "./otgrpc_testing"
)

const (
	streamLength = 5
)

type testServer struct{}

func (s *testServer) UnaryCall(ctx context.Context, in *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
	return &testpb.SimpleResponse{in.Payload}, nil
}

func (s *testServer) StreamingOutputCall(in *testpb.SimpleRequest, stream testpb.TestService_StreamingOutputCallServer) error {
	for i := 0; i < streamLength; i++ {
		if err := stream.Send(&testpb.SimpleResponse{in.Payload}); err != nil {
			return err
		}
	}
	return nil
}

func (s *testServer) StreamingInputCall(stream testpb.TestService_StreamingInputCallServer) error {
	sum := int32(0)
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += in.Payload
	}
	return stream.SendAndClose(&testpb.SimpleResponse{sum})
}

func (s *testServer) StreamingBidirectionalCall(stream testpb.TestService_StreamingBidirectionalCallServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err = stream.Send(&testpb.SimpleResponse{in.Payload}); err != nil {
			return err
		}
	}
}

type spanContext struct {
}

type span struct {
}

func (s *span) Finish() {
}

func (s *span) FinishWithOptions(opts opentracing.FinishOptions) {
}

type tracer struct {
}

type env struct {
	unaryClientInt  grpc.UnaryClientInterceptor
	streamClientInt grpc.StreamClientInterceptor
	unaryServerInt  grpc.UnaryServerInterceptor
	streamServerInt grpc.StreamServerInterceptor
}

type test struct {
	t   *testing.T
	e   env
	srv *grpc.Server
	cc  *grpc.ClientConn
	c   testpb.TestServiceClient
}

func newTest(t *testing.T, e env) *test {
	te := &test{
		t: t,
		e: e,
	}

	// Set up the server.
	sOpts := []grpc.ServerOption{}
	if e.unaryServerInt != nil {
		sOpts = append(sOpts, grpc.UnaryInterceptor(e.unaryServerInt))
	}
	if e.streamServerInt != nil {
		sOpts = append(sOpts, grpc.StreamInterceptor(e.streamServerInt))
	}
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		te.t.Fatalf("Failed to listen: %v", err)
	}
	te.srv = grpc.NewServer(sOpts...)
	testpb.RegisterTestServiceServer(te.srv, &testServer{})
	go te.srv.Serve(lis)

	// Set up a connection to the server.
	cOpts := []grpc.DialOption{grpc.WithInsecure()}
	if e.unaryClientInt != nil {
		cOpts = append(cOpts, grpc.WithUnaryInterceptor(e.unaryClientInt))
	}
	if e.streamClientInt != nil {
		cOpts = append(cOpts, grpc.WithStreamInterceptor(e.streamClientInt))
	}
	_, port, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		te.t.Fatalf("Failed to parse listener address: %v", err)
	}
	srvAddr := "localhost:" + port
	te.cc, err = grpc.Dial(srvAddr, cOpts...)
	if err != nil {
		te.t.Fatalf("Dial(%q) = %v", srvAddr, err)
	}
	te.c = testpb.NewTestServiceClient(te.cc)
	return te
}

func (te *test) tearDown() {
	te.cc.Close()
}

func TestUnaryOpenTracing(t *testing.T) {
	tracer := mocktracer.New()
	e := env{
		unaryClientInt: otgrpc.OpenTracingClientInterceptor(tracer),
		unaryServerInt: otgrpc.OpenTracingServerInterceptor(tracer),
	}
	te := newTest(t, e)
	defer te.tearDown()

	resp, err := te.c.UnaryCall(context.Background(), &testpb.SimpleRequest{0})
	if err != nil {
		t.Fatalf("Failed UnaryCall: %v", err)
	}
	assert.Equal(t, int32(0), resp.Payload)

	spans := tracer.FinishedSpans()
	assert.Equal(t, 2, len(spans))
	parent := spans[1]
	child := spans[0]
	assert.Equal(t, child.ParentID, parent.Context().(mocktracer.MockSpanContext).SpanID)
}
