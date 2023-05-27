// Obtaining a list of all warehouses and the amount of goods in them
package domain

import (
	"context"
)

// Describe information about the amount of goods in stock
type Stock struct {
	WarehouseID int64
	Count       uint64
}

// Get a list of all warehouses and the amount of goods in them
func (s *Service) GetStocks(ctx context.Context, sku uint32) ([]Stock, error) {
	return []Stock{
		{WarehouseID: 1, Count: 200},
		{WarehouseID: 2, Count: 50},
		{WarehouseID: 3, Count: 4},
		{WarehouseID: 4, Count: 71},
	}, nil
}
