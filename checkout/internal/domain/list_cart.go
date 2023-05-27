// Obtaining a list of products from the user's cart
package domain

import (
	"context"
	"errors"
	"sync"
)

// Describe the user's cart
type UserCart struct {
	Items      []Item
	TotalPrice uint32
}

var (
	ErrGetProductsInfo = errors.New("can not get all products")
)

// Get a list of items in the user's cart
func (s *Service) GetListCartWithTotalPrice(ctx context.Context, user int64) (*UserCart, error) {

	// Mock getting a list of items from the user cart
	userCart := Cart{
		OrderID: 1,
		Items: []CartItem{
			{SKU: 773297411, Count: 5},
		},
	}

	result := make([]Item, 0, len(userCart.Items))

	var wg sync.WaitGroup
	wg.Add(len(userCart.Items))
	var mu sync.Mutex

	for i := 0; i < len(userCart.Items); i++ {
		go func(index int) {
			defer wg.Done()
			item, err := s.productChecker.GetProductBySKU(ctx, userCart.Items[index].SKU)
			if err != nil {
				return
			}
			mu.Lock()
			result = append(result, Item{
				SKU:   userCart.Items[index].SKU,
				Count: userCart.Items[index].Count,
				Name:  item.Name,
				Price: item.Price,
			})
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	if len(userCart.Items) != len(result) {
		return nil, ErrGetProductsInfo
	}

	var totalPrice uint32 = 0
	for _, v := range result {
		totalPrice += uint32(v.Count) * v.Price
	}

	return &UserCart{Items: result, TotalPrice: totalPrice}, nil
}
