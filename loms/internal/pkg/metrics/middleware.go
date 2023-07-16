package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	h, err := handler(ctx, req)

	code := status.Code(err)

	prometheusMetrics.histogram.WithLabelValues(code.String(), info.FullMethod).Observe(time.Since(start).Seconds())
	prometheusMetrics.counter.With(prometheus.Labels{"type": "router"}).Inc()
	return h, err
}
