// DeleteFromCart
package cart

import (
	"context"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/cart_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteFromCart controller
func (s *Server) DeleteFromCart(ctx context.Context, req *cart_v1.DeleteFromCartRequest) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	err = s.service.DeleteFromCart(ctx, model.UserID(req.GetUser()), model.SKU(req.GetSku()), uint16(req.GetCount()))
	return &emptypb.Empty{}, err
}
