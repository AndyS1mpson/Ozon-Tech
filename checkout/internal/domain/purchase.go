// Order payment function
package domain

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

// Create a custom order
func (s *Service) Purchase(ctx context.Context, user model.UserID) (model.OrderID, error) {

	cart, err := s.cart.GetCartByUserID(ctx, user)
	if err != nil {
		return 0, errors.Wrap(err, "user have empty cart")
	}

	cartItems, err := s.cart.ListCart(ctx, cart)
	if err != nil {
		return 0, errors.Wrap(err, "get user cart items")
	}

	orderID, err := s.lomsChecker.CreateOrder(ctx, user, cartItems)
	if err != nil {
		return 0, errors.Wrap(err, "purchase order")
	}

	return orderID, nil
}
