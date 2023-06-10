// Removing an item from a user's cart
package domain

import (
	"context"
	"fmt"
)

// Delete items from the user's cart
func (s *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	err := s.cart.DeleteFromCart(ctx, user, sku, count)
	if err != nil {
		return fmt.Errorf("can not delete item: %w", err)
	}

	return nil
}
