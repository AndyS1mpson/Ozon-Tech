// Handler to cancel a user order
package cancelorder

import (
	"context"
	"route256/loms/internal/domain"
)

// Request handler
type Handler struct {
	Service *domain.Service
}

// Create New Handler instance
func New(service *domain.Service) *Handler{
	return &Handler{Service: service}
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
	err := h.Service.CancelOrder(ctx, req.OrderID)
	return Response{}, err
}
