// purchase handler
package purchase

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
	User int64 `json:"user" yaml:"user"`
}

// Describe the service response fields
type Response struct {
	OrderID int64 `json:"orderID"`
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

// Purchase request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	order, err := h.Service.Purchase(ctx, req.User)
	if err != nil {
		return Response{}, err
	}
	return Response{OrderID: order.ID}, nil
}
