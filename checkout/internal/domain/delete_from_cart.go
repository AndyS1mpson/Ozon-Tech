// Removing an item from a user's cart
package domain

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

// Delete items from the user's cart
func (s *Service) DeleteFromCart(ctx context.Context, user model.UserID, sku model.SKU, count uint16) error {
	cartID, err := s.cart.GetCartByUserID(ctx, user)
	if err != nil {
		return errors.Wrap(err, "can not get user cart")
	}

	err = s.cart.DeleteFromCart(ctx, cartID, sku, count)
	if err != nil {
		return errors.Wrap(err, "can not delete item")
	}

	return nil
}
