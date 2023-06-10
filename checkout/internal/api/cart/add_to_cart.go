// AddToCart
package cart

import (
	"context"
	"route256/checkout/pkg/cart_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AddToCart controller
func (s *Server) AddToCart(ctx context.Context, req *cart_v1.AddToCartRequest) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	err = s.service.AddToCart(
		ctx,
		req.GetUser(),
		req.GetSku(),
		uint16(req.GetCount()),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
