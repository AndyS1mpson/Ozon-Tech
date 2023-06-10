// Obtaining a list of all warehouses and the amount of goods in them
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
)

// Get a list of all warehouses and the amount of goods in them
func (s *Service) Stocks(ctx context.Context, sku model.SKU) ([]model.Stock, error) {
	stocks, err := s.stock.GetAvailableStocks(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("get stocks: %s", err)
	}

	return stocks, nil
}
