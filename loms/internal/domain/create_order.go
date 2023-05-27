// Creating a custom order
package domain

import (
	"context"
)

// Describe user good
type Item struct {
	SKU   uint32
	Count uint16
}

// Describes the user's order
type Order struct {
	User  int64
	Items []Item
}

// Create user order
func (s *Service) CreateOrder(ctx context.Context, order Order) (int64, error) {
	return 1, nil
}
