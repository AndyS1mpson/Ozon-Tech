// Client to interact with the external service loms
package loms

import (
	"context"

	"route256/checkout/internal/model"
	"route256/checkout/internal/pkg/tracer"
	"route256/checkout/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "clients/loms/get_stocks_by_sku")
	defer span.Finish()

	requestStocks := &loms_v1.StocksRequest{
		Sku: sku,
	}

	// Connect to loams service
	con, err := grpc.Dial(c.lomsAddress, grpc.WithInsecure())
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "can not connect to server loms"))
	}
	defer con.Close()

	// Create client for loms service
	lomsClient := loms_v1.NewLomsClient(con)

	// Do request
	resp, err := lomsClient.Stocks(ctx, requestStocks)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "send request error"))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "clients/loms/create_order")
	defer span.Finish()

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
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "can not connect to server loms"))
	}
	defer con.Close()

	// Create client for loms service
	lomsClient := loms_v1.NewLomsClient(con)

	// Do request
	resp, err := lomsClient.CreateOrder(ctx, requestPurchase)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "send request error"))
	}

	return model.OrderID(resp.GetOrderID()), nil
}
