// Adding an item to a user's cart
package domain

import (
	"context"
	"github.com/pkg/errors"
	"route256/checkout/internal/model"
)

var (
	ErrStockInsufficient = errors.New("stock insufficient")
)

// Add items to the user's cart if they are available
func (s *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := s.lomsChecker.GetStocksBySKU(ctx, sku)
	if err != nil {
		return errors.Wrap(err, "get stocks")
	}
	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			// Get user's cart
			cart, err := s.cart.GetCartByUserID(ctx, model.UserID(user))
			if err != nil {
				cart, err = s.cart.CreateCart(ctx, model.UserID(user))
				if err != nil {
					return errors.Wrap(err, "create cart")
				}
			}
			err = s.cart.UpdateOrAddToCart(ctx, cart, model.SKU(sku), count)
			if err != nil {
				return errors.Wrap(err, "could not add item to cart")
			}
			return nil
		}
	}

	return ErrStockInsufficient
}
