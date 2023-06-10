// Cancel order
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
)

// Cancel the order and removes the reserve from all items
func (s *Service) CancelOrder(ctx context.Context, orderID model.OrderID) error {
	order, err := s.order.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("can not get order: %w", err)
	}

	for _, item := range order.Items {
		err := s.stock.Unreserve(ctx, orderID, model.SKU(item.SKU))
		if err != nil {
			return fmt.Errorf("can not unreserve item: %w", err)
		}
	}

	err = s.order.CancelOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("try to cancel order: %s", err)
	}

	return nil
}
