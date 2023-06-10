// Get order information
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
)

// Get information about the user's order
func (s *Service) ListOrder(ctx context.Context, orderID model.OrderID) (model.OrderWithStatus, error) {
	order, err := s.order.ListOrder(ctx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, fmt.Errorf("get order info: %s", err)
	}

	return order, nil
}
