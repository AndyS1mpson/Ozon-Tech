// Obtaining a list of all warehouses and the amount of goods in them
package domain

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

// Get a list of all warehouses and the amount of goods in them
func (s *Service) Stocks(ctx context.Context, sku model.SKU) ([]model.Stock, error) {
	stocks, err := s.stock.GetAvailableStocks(ctx, sku)
	if err != nil {
		return nil, errors.Wrap(err, "get stocks")
	}

	return stocks, nil
}
