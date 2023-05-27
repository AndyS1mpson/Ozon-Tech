// listcart handler
package listcart

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
)

// Request handler
type Handler struct {
	Service *domain.Service
}

// Describe fields from the request body
type Request struct {
	User int64 `json:"user"`
}

// Describe the service response fields
type Response struct {
	Items      []ProductItem `json:"items"`
	TotalPrice uint32        `json:"totalPrice"`
}

// User product card item
type ProductItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSKUNotFound  = errors.New("sku not found")
)

// Validate data from request
func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	userCart, err := h.Service.GetListCartWithTotalPrice(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	result := make([]ProductItem, 0, len(userCart.Items))

	for _, v := range userCart.Items {
		result = append(result, ProductItem{
			SKU:   v.SKU,
			Count: v.Count,
			Name:  v.Name,
			Price: v.Price,
		})
	}

	return Response{Items: result, TotalPrice: userCart.TotalPrice}, err
}
