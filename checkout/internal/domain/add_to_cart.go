// Adding an item to a user's cart
package domain

import (
	"context"
	"errors"
	"fmt"
	"route256/checkout/internal/model"
)

var (
	ErrStockInsufficient = errors.New("stock insufficient")
)

// Add items to the user's cart if they are available
func (s *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := s.lomsChecker.GetStocksBySKU(ctx, sku)
	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
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
					return fmt.Errorf("create cart: %w", err)
				}
			}
			err = s.cart.UpdateOrAddToCart(ctx, int64(cart), sku, count)
			if err != nil {
				return fmt.Errorf("could not add item to cart: %w", err)
			}
			return nil
		}
	}

	return ErrStockInsufficient
}
