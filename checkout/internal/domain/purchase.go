// Order payment function
package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/model"
)

// Create a custom order
func (s *Service) Purchase(ctx context.Context, user model.UserID) (model.OrderID, error) {

	cart, err := s.cart.GetCartByUserID(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("user have empty cart: %w", err)
	}

	cartItems, err := s.cart.ListCart(ctx, cart)
	if err != nil {
		return 0, fmt.Errorf("get user cart items: %w", err)
	}

	orderID, err := s.lomsChecker.CreateOrder(ctx, user, cartItems)
	if err != nil {
		return 0, fmt.Errorf("purchase order: %w", err)
	}

	return orderID, nil
}
