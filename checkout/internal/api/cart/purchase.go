// Purchase
package cart

import (
	"context"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/cart_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Purchase controller
func (s *Server) Purchase(ctx context.Context, req *cart_v1.PurchaseRequest) (*cart_v1.PurchaseResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	orderId, err := s.service.Purchase(ctx, model.UserID(req.User))
	if err != nil {
		return &cart_v1.PurchaseResponse{}, status.Errorf(codes.Internal, err.Error())
	}
	return &cart_v1.PurchaseResponse{OrderID: int64(orderId)}, nil
}
