// Handler to get order info
package listorder

import (
	"context"
	"route256/loms/internal/domain"
)

// Request handler
type Handler struct {
	Service *domain.Service
}

// Describe user good
type Good struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

// Describe request body
type Request struct {
	OrderID int64 `json:"orderID"`
}

// Describe response body
type Response struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Good `json:"items"`
}

// Request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	orderInfo, err := h.Service.ListOrder(ctx, req.OrderID)
	if err != nil {
		return Response{}, err
	}

	goods := make([]Good, 0, len(orderInfo.Items))

	for _, v := range orderInfo.Items {
		goods = append(goods, Good{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}

	return Response{Status: orderInfo.Status, User: orderInfo.User, Items: goods}, err
}
