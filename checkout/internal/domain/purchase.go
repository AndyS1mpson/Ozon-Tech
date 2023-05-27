// Order payment function
package domain

import (
	"context"
	"fmt"
)

// Create a custom order
func (s *Service) Purchase(ctx context.Context, user int64) (OrderID, error) {

	// Mock getting a list of items from the user cart
	cartItems := []CartItem{
		{SKU: 1, Count: 5},
		{SKU: 2, Count: 2},
		{SKU: 3, Count: 10},
	}

	order, err := s.lomsChecker.CreateOrder(ctx, user, cartItems)
	if err != nil {
		return OrderID{}, fmt.Errorf("purchase order: %w", err)
	}

	return order, nil
}
