// Get order information
package domain

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

// Get information about the user's order
func (s *Service) ListOrder(ctx context.Context, orderID model.OrderID) (model.OrderWithStatus, error) {
	order, err := s.order.ListOrder(ctx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, errors.Wrap(err, "get order info")
	}

	return order, nil
}
