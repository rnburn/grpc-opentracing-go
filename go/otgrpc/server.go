package otgrpc

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// OpenTracingServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
//
// For example:
//
//     s := grpc.NewServer(
//         ...,  // (existing ServerOptions)
//         grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
//
// All gRPC server spans will look for an OpenTracing SpanContext in the gRPC
// metadata; if found, the server span will act as the ChildOf that RPC
// SpanContext.
//
// Root or not, the server Span will be embedded in the context.Context for the
// application-specific gRPC handler(s) to access.
func OpenTracingServerInterceptor(tracer opentracing.Tracer, optFuncs ...Option) grpc.UnaryServerInterceptor {
	otgrpcOpts := newOptions()
	otgrpcOpts.apply(optFuncs...)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		spanContext, err := extractSpanContext(ctx, tracer)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			// TODO: establish some sort of error reporting mechanism here. We
			// don't know where to put such an error and must rely on Tracer
			// implementations to do something appropriate for the time being.
		}
		if otgrpcOpts.inclusionFunc != nil &&
			!otgrpcOpts.inclusionFunc(spanContext, info.FullMethod, req, nil) {
			return handler(ctx, req)
		}
		serverSpan := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			gRPCComponentTag,
		)
		defer serverSpan.Finish()

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		if otgrpcOpts.logPayloads {
			serverSpan.LogFields(log.Object("gRPC request", req))
		}
		resp, err = handler(ctx, req)
		if err == nil {
			if otgrpcOpts.logPayloads {
				serverSpan.LogFields(log.Object("gRPC response", resp))
			}
		} else {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "gRPC error"), log.Error(err))
		}
		if otgrpcOpts.decorator != nil {
			otgrpcOpts.decorator(serverSpan, info.FullMethod, req, resp, err)
		}
		return resp, err
	}
}

func OpenTracingStreamServerInterceptor(tracer opentracing.Tracer, optFuncs ...Option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		spanContext, err := extractSpanContext(ss.Context(), tracer)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			// TODO: establish some sort of error reporting mechanism here. We
			// don't know where to put such an error and must rely on Tracer
			// implementations to do something appropriate for the time being.
		}
		serverSpan := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			gRPCComponentTag,
		)
		defer serverSpan.Finish()
		ss = &openTracingServerStream{ss, opentracing.ContextWithSpan(ss.Context(), serverSpan)}
		err = handler(srv, ss)
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "gRPC error"), log.Error(err))
		}
		return err
	}
}

type openTracingServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (stream *openTracingServerStream) Context() context.Context {
	return stream.ctx
}

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
}
