package repository

import (
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/schema"
)

// Convert db object to domain model object
func ToCartItem(cartItem schema.CartItem) model.CartItem {
	return model.CartItem{
		SKU:   cartItem.SKU,
		Count: cartItem.Count,
	}
}

// Convert user cart items from db to domain model objects
func ToCartItems(cartItems []schema.CartItem) []model.CartItem {
	result := make([]model.CartItem, len(cartItems))
	for i, cartItem := range cartItems {
		result[i] = ToCartItem(cartItem)
	}
	return result
}
