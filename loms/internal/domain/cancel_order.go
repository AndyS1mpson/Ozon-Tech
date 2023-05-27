// Cancel order
package domain

import (
	"context"
)

// Cancel the order and removes the reserve from all items
func (s *Service) CancelOrder(ctx context.Context, orderID int64) error {
	return nil
}
