// Client to interact with the external product service
package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"route256/checkout/internal/domain"
	"route256/libs/clientwrapper"
)

const (
	productsPath = "get_product"
)

// Describe fields of the body of the request to the product service
type ProductsRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

// Describe the response body of the product service
type ProductsResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

// Implement interaction with the product service
type Client struct {
	pathProducts string
	token        string
}

// Creates a new client instance
func New(clientUrl string, token string) *Client {
	productUrl, _ := url.JoinPath(clientUrl, productsPath)
	return &Client{pathProducts: productUrl, token: token}
}

// Show a list of products in the user's cart
func (c *Client) GetProductBySKU(ctx context.Context, sku uint32) (domain.ProductInfo, error) {
	requestProduct := ProductsRequest{Token: c.token, SKU: sku}

	httpResponse, err := clientwrapper.DoRequest(ctx, requestProduct, c.pathProducts, "POST")
	if err != nil {
		return domain.ProductInfo{}, err
	}

	responseProduct := ProductsResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responseProduct)
	if err != nil {
		return domain.ProductInfo{}, fmt.Errorf("decode stock request: %w", err)
	}

	return domain.ProductInfo{
		Name:  responseProduct.Name,
		Price: responseProduct.Price,
	}, nil
}
