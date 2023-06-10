// OrderPayed
package loms

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// OrderPayed controller
func (s *Server) OrderPayed(ctx context.Context, req *loms_v1.OrderPayedRequest) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	err = s.service.OrderPayed(ctx, model.OrderID(req.GetOrderID()))
	return &emptypb.Empty{}, err
}
