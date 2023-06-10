// ListOrder
package loms

import (
	"context"
	"route256/loms/internal/converter/server"
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListOrder controller
func (s *Server) ListOrder(ctx context.Context, req *loms_v1.ListOrderRequest) (*loms_v1.ListOrderResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	orderInfo, err := s.service.ListOrder(ctx, model.OrderID(req.GetOrderID()))
	if err != nil {
		return &loms_v1.ListOrderResponse{}, status.Errorf(codes.Internal, err.Error())
	}
	res := server.ListOrderToResp(orderInfo)
	return &res, nil
}
