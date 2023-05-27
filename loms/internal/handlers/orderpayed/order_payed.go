// Handler to pay for the order
package orderpayed

import (
	"context"
	"route256/loms/internal/domain"
)

// Request handler
type Handler struct {
	Service *domain.Service
}

// Describe request body
type Request struct {
	OrderID int64 `json:"orderID"`
}

// Describe response body
type Response struct {
}

// Request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	err := h.Service.OrderPayed(ctx, req.OrderID)

	return Response{}, err
}
