// Converters for transferring objects between layers
package server

import (
	"route256/checkout/internal/model"
	"route256/checkout/pkg/cart_v1"
)

// Convert Good to response object
func CartItemToRe(item model.Good) *cart_v1.CartGoodInfo {
	return &cart_v1.CartGoodInfo{
		Sku:   item.SKU,
		Count: uint32(item.Count),
		Name:  item.Name,
		Price: item.Price,
	}
}

// Convert UserCartWithTotal to response object
func ListCartToRe(req model.UserCartWithTotal) *cart_v1.ListCartResponse {
	items := []*cart_v1.CartGoodInfo{}
	for _, item := range req.Items {
		items = append(items, CartItemToRe(item))
	}
	return &cart_v1.ListCartResponse{
		Items:      items,
		TotalPrice: req.TotalPrice,
	}
}
