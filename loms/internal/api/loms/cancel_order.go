// CancelOrder
package loms

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CancelOrder controller
func (s *Server) CancelOrder(ctx context.Context, req *loms_v1.CancelOrderRequest) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	err = s.service.CancelOrder(ctx, model.OrderID(req.GetOrderID()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
