// Description of things common to the domain layer
package domain

import (
	"context"
	"route256/loms/internal/model"
)

// Describe repository for working with orders
type OrderRepository interface {
	GetOrder(ctx context.Context, id model.OrderID) (*model.Order, error)
	CreateOrder(ctx context.Context,  order model.Order) (model.OrderID, error)
	ListOrder(ctx context.Context, orderID model.OrderID) (model.OrderWithStatus, error)
	PayOrder(ctx context.Context, orderID model.OrderID) error
	CancelOrder(ctx context.Context, orderID model.OrderID) error
	AwaitPaymentOrder(ctx context.Context, orderID model.OrderID) error
	FailOrder(ctx context.Context, orderID model.OrderID) error
}

// Describe repository for working with stocks
type StockRepository interface {
	GetAvailableStocks(ctx context.Context, sku model.SKU) ([]model.Stock, error)
	Reserve(ctx context.Context, orderID model.OrderID, sku model.SKU, stock model.Stock) error
	Unreserve(ctx context.Context, orderID model.OrderID, sku model.SKU) error
	WriteOffOrderItems(ctx context.Context, orderID model.OrderID) ([]model.Stock, error)
}

// Provide access to the business logic of the service
type Service struct {
	order OrderRepository
	stock StockRepository
}

// Create a new Service instance
func New(order OrderRepository, stock StockRepository) *Service {
	return &Service{
		order: order,
		stock: stock,
	}
}
