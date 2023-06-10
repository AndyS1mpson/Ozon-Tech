// ListToCart
package cart

import (
	"context"
	"route256/checkout/internal/converter/server"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/cart_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListToCart controller
func (s *Server) ListCart(ctx context.Context, req *cart_v1.ListCartRequest) (*cart_v1.ListCartResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	cart, err := s.service.ListCart(ctx, model.UserID(req.GetUser()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res, err := server.ListCartToRe(cart)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return res, nil
}
