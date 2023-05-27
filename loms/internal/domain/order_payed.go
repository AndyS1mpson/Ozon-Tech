// Changing order status to paid
package domain

import (
	"context"
)

// Mark the order as paid
func (s *Service) OrderPayed(ctx context.Context, orderID int64) error {
	return nil
}
