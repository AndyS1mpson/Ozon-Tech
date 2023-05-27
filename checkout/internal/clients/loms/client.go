// Client to interact with the external service loms
package loms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"route256/checkout/internal/domain"
	"route256/libs/clientwrapper"
)

const (
	stocksPath   = "stocks"
	purchasePath = "createOrder"
)

// Describe fields of the body of the stocks request to the service loms
type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

// Describe the bodies of the request from the stocks endpoint
// of the loms service
type StocksResponse struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

// Implement interaction with the loms service
type Client struct {
	pathStock    string
	pathPurchase string
}

// Creates a new client instance
func New(clientUrl string) *Client {
	stockUrl, _ := url.JoinPath(clientUrl, stocksPath)
	purchaseUrl, _ := url.JoinPath(clientUrl, purchasePath)
	return &Client{pathStock: stockUrl, pathPurchase: purchaseUrl}
}

// Get the quantity of goods from all warehouses from the service loms
func (c *Client) GetStocksBySKU(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	requestStocks := StocksRequest{SKU: sku}

	httpResponse, err := clientwrapper.DoRequest(ctx, requestStocks, c.pathStock, "GET")
	if err != nil {
		return []domain.Stock{}, err
	}

	responseStocks := StocksResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseStocks)
	if err != nil {
		return nil, fmt.Errorf("decode stock request: %w", err)
	}

	result := make([]domain.Stock, 0, len(responseStocks.Stocks))
	for _, v := range responseStocks.Stocks {
		result = append(result, domain.Stock{
			WarehouseID: uint64(v.WarehouseID),
			Count:       v.Count,
		})
	}

	return result, nil
}

// Describe user good
type PurchaseItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

// Describe fields of the body of the purchase request to the service loms
type PurchaseRequest struct {
	User  int64          `json:"user"`
	Items []PurchaseItem `json:"items"`
}

// Describe the bodies of the request from the purchase endpoint
// of the loms service
type PurchaseResponse struct {
	OrderID int64 `json:"orderID"`
}

// Create user order
func (c *Client) CreateOrder(ctx context.Context, user int64, userGoods []domain.CartItem) (domain.OrderID, error) {
	items := make([]PurchaseItem, 0, len(userGoods))
	for _, v := range userGoods {
		items = append(items, PurchaseItem{
			SKU:   v.SKU,
			Count: v.Count,
		})
	}
	requestPurchase := PurchaseRequest{User: user, Items: items}

	httpResponse, err := clientwrapper.DoRequest(ctx, requestPurchase, c.pathPurchase, "GET")
	if err != nil {
		return domain.OrderID{}, err
	}

	responseOrder := PurchaseResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseOrder)
	if err != nil {
		return domain.OrderID{}, fmt.Errorf("decode stock request: %w", err)
	}

	return domain.OrderID{ID: responseOrder.OrderID}, nil
}
