// Get order information
package domain

import (
	"context"
)

// Get the user's order status
type OrderStatus struct {
	Status string
	User   int64
	Items  []Item
}

// Get information about the user's order
func (s *Service) ListOrder(ctx context.Context, orderID int64) (OrderStatus, error) {
	return OrderStatus{
		Status: "payed",
		User:   1,
		Items: []Item{
			{SKU: 1, Count: 5},
			{SKU: 2, Count: 5},
		},
	}, nil
}
