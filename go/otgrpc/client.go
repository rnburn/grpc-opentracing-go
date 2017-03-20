package otgrpc

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
)

// OpenTracingClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
//
// For example:
//
//     conn, err := grpc.Dial(
//         address,
//         ...,  // (existing DialOptions)
//         grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
//
// All gRPC client spans will inject the OpenTracing SpanContext into the gRPC
// metadata; they will also look in the context.Context for an active
// in-process parent Span and establish a ChildOf reference if such a parent
// Span could be found.
func OpenTracingClientInterceptor(tracer opentracing.Tracer, optFuncs ...Option) grpc.UnaryClientInterceptor {
	otgrpcOpts := newOptions()
	otgrpcOpts.apply(optFuncs...)
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		if otgrpcOpts.inclusionFunc != nil &&
			!otgrpcOpts.inclusionFunc(parentCtx, method, req, resp) {
			return invoker(ctx, method, req, resp, cc, opts...)
		}
		clientSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			gRPCComponentTag,
		)
		defer clientSpan.Finish()
		ctx = injectSpanContext(ctx, tracer, clientSpan)
		if otgrpcOpts.logPayloads {
			clientSpan.LogFields(log.Object("gRPC request", req))
		}
		err = invoker(ctx, method, req, resp, cc, opts...)
		if err == nil {
			if otgrpcOpts.logPayloads {
				clientSpan.LogFields(log.Object("gRPC response", resp))
			}
		} else {
			clientSpan.LogFields(log.String("event", "gRPC error"), log.Error(err))
			ext.Error.Set(clientSpan, true)
		}
		if otgrpcOpts.decorator != nil {
			otgrpcOpts.decorator(clientSpan, method, req, resp, err)
		}
		return err
	}
}

// OpenTracingStreamClientInterceptor returns a grpc.StreamClientInterceptor suitable
// for use in a grpc.Dial call.
//
// For example:
//
//     conn, err := grpc.Dial(
//         address,
//         ...,  // (existing DialOptions)
//         grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)))
//
// All gRPC client spans will inject the OpenTracing SpanContext into the gRPC
// metadata; they will also look in the context.Context for an active
// in-process parent Span and establish a ChildOf reference if such a parent
// Span could be found.
func OpenTracingStreamClientInterceptor(tracer opentracing.Tracer, optFuncs ...Option) grpc.StreamClientInterceptor {
	otgrpcOpts := newOptions()
	otgrpcOpts.apply(optFuncs...)
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		var err error
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		if otgrpcOpts.inclusionFunc != nil &&
			!otgrpcOpts.inclusionFunc(parentCtx, method, nil, nil) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		clientSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			gRPCComponentTag,
		)
		ctx = injectSpanContext(ctx, tracer, clientSpan)
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			clientSpan.LogFields(log.String("event", "gRPC error"), log.Error(err))
			ext.Error.Set(clientSpan, true)
			clientSpan.Finish()
			return cs, err
		}
		return newOpenTracingClientStream(cs, method, desc, tracer, clientSpan, otgrpcOpts), nil
	}
}

func newOpenTracingClientStream(cs grpc.ClientStream, method string, desc *grpc.StreamDesc, tracer opentracing.Tracer,
	clientSpan opentracing.Span, otgrpcOpts *options) grpc.ClientStream {
	finishChan := make(chan struct{})
	lock := new(sync.Mutex)
	go func() {
		select {
		case <-finishChan:
			// The client span is being finished by another code path; hence, no
			// action is necessary.
		case <-cs.Context().Done():
			finishStreamSpan(lock, finishChan, clientSpan, method, otgrpcOpts.decorator, cs.Context().Err())
		}
	}()
	otcs := &openTracingClientStream{
		ClientStream: cs,
		finishChan:   finishChan,
		method:       method,
		desc:         desc,
		otgrpcOpts:   otgrpcOpts,
		lock:         lock,
		clientSpan:   clientSpan,
	}
	return otcs
}

type openTracingClientStream struct {
	grpc.ClientStream
	finishChan chan struct{}
	method     string
	desc       *grpc.StreamDesc
	otgrpcOpts *options
	lock       *sync.Mutex
	clientSpan opentracing.Span
}

func (cs *openTracingClientStream) Header() (metadata.MD, error) {
	md, err := cs.ClientStream.Header()
	if err != nil {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, err)
	}
	return md, err
}

func (cs *openTracingClientStream) SendMsg(m interface{}) error {
	if cs.otgrpcOpts.logPayloads {
		cs.logMsg("gRPC request", m)
	}
	err := cs.ClientStream.SendMsg(m)
	if err != nil {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, err)
	}
	return err
}

func (cs *openTracingClientStream) RecvMsg(m interface{}) error {
	err := cs.ClientStream.RecvMsg(m)
	if err == io.EOF {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, nil)
		return err
	} else if err != nil {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, err)
		return err
	}
	if cs.otgrpcOpts.logPayloads {
		cs.logMsg("gRPC response", m)
	}
	if !cs.desc.ServerStreams {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, nil)
	}
	return err
}

func (cs *openTracingClientStream) CloseSend() error {
	err := cs.ClientStream.CloseSend()
	if err != nil {
		finishStreamSpan(cs.lock, cs.finishChan, cs.clientSpan, cs.method, cs.otgrpcOpts.decorator, err)
	}
	return err
}

func (cs *openTracingClientStream) logMsg(name string, m interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	select {
	case <-cs.finishChan:
		// The client span is closed or being closed, so we can't log anything.
		return
	default:
	}
	cs.clientSpan.LogFields(log.Object(name, m))
}

func injectSpanContext(ctx context.Context, tracer opentracing.Tracer, clientSpan opentracing.Span) context.Context {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := metadataReaderWriter{md}
	err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdWriter)
	// We have no better place to record an error than the Span itself :-/
	if err != nil {
		clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
	}
	return metadata.NewContext(ctx, md)
}

func finishStreamSpan(lock *sync.Mutex, finishChan chan struct{}, clientSpan opentracing.Span,
	method string, decorator SpanDecoratorFunc, err error) {
	lock.Lock()
	defer lock.Unlock()
	select {
	case <-finishChan:
		// The client span is either already finished or being finished, so we have
		// nothing to do.
		return
	default:
	}
	close(finishChan)
	if err != nil {
		clientSpan.LogFields(log.String("event", "gRPC error"), log.Error(err))
		ext.Error.Set(clientSpan, true)
	}
	if decorator != nil {
		decorator(clientSpan, method, nil, nil, err)
	}
	clientSpan.Finish()
}
