// Handler to create a user order
package createorder

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

// Describe user good
type Good struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

// Describe request body
type Request struct {
	User  int64  `json:"user"`
	Items []Good `json:"items"`
}

// Describe response body
type Response struct {
	OrderID int64 `json:"orderID"`
}

// Request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	goods := make([]domain.Item, 0, len(req.Items))

	for _, v := range req.Items {
		goods = append(goods, domain.Item{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}
	orderID, err := h.Service.CreateOrder(ctx, domain.Order{User: req.User, Items: goods})

	return Response{OrderID: orderID}, err
}
