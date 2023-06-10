// Obtaining a list of products from the user's cart
package domain

import (
	"context"
	"errors"
	"fmt"
	"route256/checkout/internal/model"
	"sync"
)

var (
	ErrGetProductsInfo = errors.New("can not get all products")
)

// Get a list of items in the user's cart
func (s *Service) ListCart(ctx context.Context, user model.UserID) (model.UserCartWithTotal, error) {
	cart, err := s.cart.GetCartByUserID(ctx, user)
	if err != nil {
		cart, err = s.cart.CreateCart(ctx, user)
		if err != nil {
			return model.UserCartWithTotal{}, fmt.Errorf("error creating cart: %v", err)
		}
	}

	userCart, err := s.cart.ListCart(ctx, cart)
	if err != nil {
		return model.UserCartWithTotal{}, fmt.Errorf("can not get cart info: %v", err)
	}

	result := make([]model.Good, 0, len(userCart))

	var wg sync.WaitGroup
	wg.Add(len(userCart))
	var mu sync.Mutex

	for i := 0; i < len(userCart); i++ {
		go func(index int) {
			defer wg.Done()
			item, err := s.productChecker.GetProduct(ctx, userCart[index].SKU)
			if err != nil {
				return
			}
			mu.Lock()
			result = append(result, model.Good{
				SKU:   userCart[index].SKU,
				Count: userCart[index].Count,
				Name:  item.Name,
				Price: item.Price,
			})
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	if len(userCart) != len(result) {
		return model.UserCartWithTotal{}, ErrGetProductsInfo
	}

	var totalPrice uint32 = 0
	for _, v := range result {
		totalPrice += uint32(v.Count) * v.Price
	}

	return model.UserCartWithTotal{Items: result, TotalPrice: totalPrice}, nil
}
