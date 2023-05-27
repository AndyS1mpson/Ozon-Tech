// Adding an item to a user's cart
package domain

import (
	"context"
	"errors"
	"fmt"
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
			return nil
		}
	}

	return ErrStockInsufficient
}
