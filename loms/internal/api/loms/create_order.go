// CreateOrder
package loms

import (
	"context"
	"route256/loms/internal/converter/server"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateOrder controller
func (s *Server) CreateOrder(ctx context.Context, req *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {
	or, err := server.OrderFromReq(req)
	if err != nil {
		return &loms_v1.CreateOrderResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
	}
	orderID, err := s.service.CreateOrder(ctx, or)
	if err != nil {
		return &loms_v1.CreateOrderResponse{}, status.Errorf(codes.Internal, err.Error())
	}
	return &loms_v1.CreateOrderResponse{OrderID: int64(orderID)}, nil
}
