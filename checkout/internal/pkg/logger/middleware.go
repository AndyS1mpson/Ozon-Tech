package logger

import (
	"context"

	"google.golang.org/grpc"
)

// Middleware for logging requests
func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	h, err := handler(ctx, req)
	if err != nil {
		Error(info.FullMethod, "error while processing handler, err: ", err)
	}

	return h, err
}
