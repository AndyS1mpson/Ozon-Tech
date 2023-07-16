// Client to interact with the external product service
package products

import (
	"context"
	"log"
	"route256/checkout/internal/model"
	"route256/checkout/internal/pkg/ratelimit"
	"route256/checkout/internal/pkg/tracer"
	"route256/checkout/internal/pkg/workerpool"
	"route256/checkout/pkg/product_v1"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	productsPath = "get_product"
	workerCount  = 10
)

// Implement interaction with the product service
type Client struct {
	productAddress string
	token          string
	limiter        ratelimit.RateLimiter
}

// Creates a new client instance
func New(address string, token string, limiter ratelimit.RateLimiter) *Client {
	return &Client{productAddress: address, token: token, limiter: limiter}
}

// Get info about list of products from the user's cart
func (c *Client) GetProducts(ctx context.Context, goods []model.CartItem) ([]model.Good, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "clients/products/get_products")
	defer span.Finish()

	// Connect to loams service
	con, err := grpc.Dial(c.productAddress, grpc.WithInsecure())
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "can not connect to server loms"))
	}
	defer con.Close()

	result := make([]model.Good, 0, len(goods))
	// Create client for loms service
	productClient := product_v1.NewProductServiceClient(con)

	pool := workerpool.NewPool[ClientWithData, model.Good](ctx, workerCount)

	var wg sync.WaitGroup
	wg.Add(len(goods))
	for i := 0; i < len(goods); i++ {
		// Wait if the limit of requests per second has already been reached
		err := c.limiter.Acquire(ctx)
		if err != nil {
			return nil, tracer.MarkSpanWithError(ctx, err)
		}

		// Add task to pool
		taskData := ClientWithData{
			Good:  goods[i],
			Token: c.token,
			Con:   productClient,
		}
		p := pool.Exec(taskData, getProduct)

		go func(index int) {
			defer wg.Done()
			defer c.limiter.Release()

			good := <-p.Out
			if good.Err != nil {
				log.Printf("ERROR: can not get product info: %v", good.Err)
				return
			}
			result = append(result, model.Good{
				SKU:   goods[index].SKU,
				Count: goods[index].Count,
				Name:  good.Value.Name,
				Price: good.Value.Price,
			})
		}(i)
	}
	wg.Wait()
	c.limiter.Close()

	return result, nil
}

type ClientWithData struct {
	Good  model.CartItem
	Token string
	Con   product_v1.ProductServiceClient
}

// Get info for one product
func getProduct(ctx context.Context, request ClientWithData) (model.Good, error) {
	requestProduct := &product_v1.GetProductRequest{
		Token: request.Token,
		Sku:   uint32(request.Good.SKU),
	}

	// Do request
	resp, err := request.Con.GetProduct(ctx, requestProduct)
	if err != nil {
		return model.Good{}, errors.Wrap(err, "send request error")
	}
	return model.Good{
		Name:  resp.Name,
		Price: resp.Price,
		SKU:   request.Good.SKU,
		Count: request.Good.Count,
	}, nil
}
