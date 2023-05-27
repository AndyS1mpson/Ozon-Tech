// deletefromcart handler
package deletefromcart

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
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

// Describe the service response fields
type Response struct {
}

var (
	ErrUserNotFound = errors.New("user not found")
)

// Validate data from request
func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Request handler function
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	err := h.Service.DeleteFromCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, err
}
