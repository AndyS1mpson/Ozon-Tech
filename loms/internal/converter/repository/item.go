package repository

import (
	"route256/loms/internal/model"
	schema "route256/loms/internal/repository/scheme"
)

// Convert db object to domain model object
func ToOrderItem(cartItem schema.Item) model.OrderItem {
	return model.OrderItem{
		SKU:   cartItem.SKU,
		Count: cartItem.Count,
	}
}

// Convert user order items from db to domain model objects
func ToOrderItems(cartItems []schema.Item) []model.OrderItem {
	result := make([]model.OrderItem, len(cartItems))
	for i, cartItem := range cartItems {
		result[i] = ToOrderItem(cartItem)
	}
	return result
}
