// Obtaining a list of products from the user's cart
package domain

import (
	"context"
	"github.com/pkg/errors"
	"route256/checkout/internal/model"
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
			return model.UserCartWithTotal{}, errors.Wrap(err, "error creating cart")
		}
		return model.UserCartWithTotal{}, nil
	}

	userCart, err := s.cart.ListCart(ctx, cart)
	if err != nil {
		return model.UserCartWithTotal{}, errors.Wrap(err, "can not get cart info")
	}

	result, err := s.productChecker.GetProducts(ctx, userCart)
	if err != nil {
		return model.UserCartWithTotal{}, errors.Wrap(err, "can not get products info")
	}

	if len(userCart) != len(result) {
		return model.UserCartWithTotal{}, ErrGetProductsInfo
	}

	var totalPrice uint32 = 0
	for _, v := range result {
		totalPrice += uint32(v.Count) * v.Price
	}

	return model.UserCartWithTotal{Items: result, TotalPrice: totalPrice}, nil
}
