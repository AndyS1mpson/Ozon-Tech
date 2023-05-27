// Removing an item from a user's cart
package domain

import "context"

// Delete items from the user's cart
func (s *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	return nil
}
