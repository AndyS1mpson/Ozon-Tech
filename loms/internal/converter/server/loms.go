// Converters for the presentation layer
package server

import (
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"
)

// Convert order item from request to OrderItem
func OrderItemFromReq(item *loms_v1.OrderItem) (*model.OrderItem, error) {
	err := item.ValidateAll()
	if err != nil {
		return nil, err
	}
	return &model.OrderItem{
		SKU:   item.GetSku(),
		Count: uint16(item.GetCount()),
	}, nil
}

// Convert OrderItem to response object
func OrderItemToRes(item model.OrderItem) *loms_v1.OrderItem {
	return &loms_v1.OrderItem{
		Sku:   item.SKU,
		Count: uint32(item.Count),
	}
}

// Convert data from request to Order
func OrderFromReq(req *loms_v1.CreateOrderRequest) (model.Order, error) {
	items := []model.OrderItem{}

	for _, item := range req.GetItems() {
		orderItem, err := OrderItemFromReq(item)
		if err != nil {
			return model.Order{}, err
		}
		items = append(items, *orderItem)
	}

	return model.Order{
		User:  req.GetUser(),
		Items: items,
	}, nil
}

// Convert order info to response object
func ListOrderToResp(order model.OrderWithStatus) loms_v1.ListOrderResponse {
	items := []*loms_v1.OrderItem{}
	for _, or := range order.Items {
		items = append(items, OrderItemToRes(or))
	}

	return loms_v1.ListOrderResponse{
		Status: order.Status,
		User:   order.User,
		Items:  items,
	}
}

// Convert stock item to response object
func StockToRes(stock model.Stock) *loms_v1.Stock {
	return &loms_v1.Stock{
		WarehouseID: stock.WarehouseID,
		Count:       stock.Count,
	}
}

// Convert stock items to response object
func StocksToRes(stocks []model.Stock) loms_v1.StocksResponse {
	items := []*loms_v1.Stock{}
	for _, it := range stocks {
		items = append(items, StockToRes(it))
	}

	return loms_v1.StocksResponse{
		Stocks: items,
	}
}
