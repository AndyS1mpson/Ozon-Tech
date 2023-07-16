package tracer

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

// Middleware for tracing requests
func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		_ = grpc.SendHeader(ctx, map[string][]string{
			"x-trace-id": {spanContext.TraceID().String()},
		})
	}

	h, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
	}

	return h, err
}

// Mark span as returned error
func MarkSpanWithError(ctx context.Context, err error) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return err
	}

	ext.Error.Set(span, true)
	span.LogKV("error", err.Error())

	return err
}
