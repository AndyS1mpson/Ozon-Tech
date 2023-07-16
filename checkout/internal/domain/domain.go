//go:generate mockery --output ./mocks --filename loms_checker_mock.go --name LomsChecker
//go:generate mockery --output ./mocks --filename product_checker_mock.go --name ProductChecker
//go:generate mockery --output ./mocks --filename cart_repository_mock.go --name CartRepository
package domain

import (
	"context"
	"route256/checkout/internal/model"
)

// Describe methods to check the availability of goods in stock
type LomsChecker interface {
	GetStocksBySKU(ctx context.Context, sku uint32) ([]model.Stock, error)
	CreateOrder(ctx context.Context, user model.UserID, userGoods []model.CartItem) (model.OrderID, error)
}

// Describes the method of retrieving product information
type ProductChecker interface {
	GetProducts(ctx context.Context, goods []model.CartItem) ([]model.Good, error)
}

type CartRepository interface {
	CreateCart(ctx context.Context, user model.UserID) (model.UserCartID, error)
	GetCartByUserID(ctx context.Context, userID model.UserID) (model.UserCartID, error)
	UpdateOrAddToCart(ctx context.Context, cart model.UserCartID, sku model.SKU, count uint16) error
	DeleteFromCart(ctx context.Context, user model.UserCartID, sku model.SKU, count uint16) error
	ListCart(ctx context.Context, cart model.UserCartID) ([]model.CartItem, error)
}

// Provide access to the business logic of the service
type Service struct {
	lomsChecker    LomsChecker
	productChecker ProductChecker
	cart           CartRepository
}

// Create a new Service instance
func New(lomsChecker LomsChecker, productChecker ProductChecker, cart CartRepository) *Service {
	return &Service{
		lomsChecker:    lomsChecker,
		productChecker: productChecker,
		cart:           cart,
	}
}
