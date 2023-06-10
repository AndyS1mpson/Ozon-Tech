// Client to interact with the external service loms
package loms

import (
	"context"

	"fmt"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/loms_v1"

	"google.golang.org/grpc"
)

// Implement interaction with the loms service
type Client struct {
	lomsAddress string
}

// Creates a new client instance
func New(clientAddress string) *Client {
	return &Client{lomsAddress: clientAddress}
}

// Get the quantity of goods from all warehouses from the service loms
func (c *Client) GetStocksBySKU(ctx context.Context, sku uint32) ([]model.Stock, error) {
	requestStocks := &loms_v1.StocksRequest{
		Sku: sku,
	}

	// Connect to loams service
	con, err := grpc.Dial(c.lomsAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("can not connect to server loms: %v", err)
	}
	defer con.Close()

	// Create client for loms service
	lomsClient := loms_v1.NewLomsClient(con)

	// Do request
	resp, err := lomsClient.Stocks(ctx, requestStocks)
	if err != nil {
		return nil, fmt.Errorf("send request error: %v", err)
	}

	items := []model.Stock{}
	for _, item := range resp.GetStocks() {
		items = append(items, model.Stock{
			WarehouseID: uint64(item.GetWarehouseID()),
			Count:       item.GetCount(),
		})
	}

	return items, nil
}

// Create user order
func (c *Client) CreateOrder(ctx context.Context, user model.UserID, userGoods []model.CartItem) (model.OrderID, error) {
	items := make([]*loms_v1.OrderItem, 0, len(userGoods))
	for _, v := range userGoods {
		items = append(items, &loms_v1.OrderItem{
			Sku:   v.SKU,
			Count: uint32(v.Count),
		})
	}
	requestPurchase := &loms_v1.CreateOrderRequest{
		User:  int64(user),
		Items: items,
	}

	// Connect to loams service
	con, err := grpc.Dial(c.lomsAddress, grpc.WithInsecure())
	if err != nil {
		return 0, fmt.Errorf("can not connect to server loms: %v", err)
	}
	defer con.Close()

	// Create client for loms service
	lomsClient := loms_v1.NewLomsClient(con)

	// Do request
	resp, err := lomsClient.CreateOrder(ctx, requestPurchase)
	if err != nil {
		return 0, fmt.Errorf("send request error: %v", err)
	}

	return model.OrderID(resp.GetOrderID()), nil
}
