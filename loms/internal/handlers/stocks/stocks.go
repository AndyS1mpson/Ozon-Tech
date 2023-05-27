// Handler for getting information about the quantity of goods in each warehouse
package stocks

import (
	"context"
	"route256/loms/internal/domain"
)

// Request handler
type Handler struct {
	Service *domain.Service
}

// Describe the quantity of goods in stock
type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

// Describe request body
type Request struct {
	SKU uint32 `json:"sku"`
}

// Describe response body
type Response struct {
	Stocks []StockItem `json:"stocks"`
}

// Request handler
func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	stocks, err := h.Service.GetStocks(ctx, req.SKU)

	result := make([]StockItem, 0, len(stocks))

	for _, v := range stocks {
		result = append(result, StockItem{
			WarehouseID: v.WarehouseID,
			Count:       v.Count,
		})
	}
	return Response{Stocks: result}, err
}
