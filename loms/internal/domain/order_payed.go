// Changing order status to paid
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
)

// Mark the order as paid
func (s *Service) OrderPayed(ctx context.Context, orderID model.OrderID) error {
	_, err := s.stock.WriteOffOrderItems(ctx, orderID)
	if err != nil {
		return fmt.Errorf("try to write off order items from warehouses: %s", err)
	}

	err = s.order.PayOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("try to paid order: %s", err)
	}

	return nil
}
