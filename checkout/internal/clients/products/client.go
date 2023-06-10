// Client to interact with the external product service
package products

import (
	"context"
	"fmt"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/product_v1"

	"google.golang.org/grpc"
)

const (
	productsPath = "get_product"
)

// Implement interaction with the product service
type Client struct {
	productAddress string
	token          string
}

// Creates a new client instance
func New(address string, token string) *Client {
	return &Client{productAddress: address, token: token}
}

// Show a list of products in the user's cart
func (c *Client) GetProduct(ctx context.Context, sku uint32) (model.Product, error) {
	requestProduct := &product_v1.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	}

	// Connect to loams service
	con, err := grpc.Dial(c.productAddress, grpc.WithInsecure())
	if err != nil {
		return model.Product{}, fmt.Errorf("can not connect to server loms: %v", err)
	}
	defer con.Close()

	// Create client for loms service
	productClient := product_v1.NewProductServiceClient(con)
	// Do request
	resp, err := productClient.GetProduct(ctx, requestProduct)
	if err != nil {
		return model.Product{}, fmt.Errorf("send request error: %v", err)
	}
	return model.Product{
		Name:  resp.Name,
		Price: resp.Price,
	}, nil
}
