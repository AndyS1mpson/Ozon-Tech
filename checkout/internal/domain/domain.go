// Description of things common to the domain layer
package domain

import "context"

// Describe the amount of goods in stock
type Stock struct {
	WarehouseID uint64
	Count       uint64
}

// Product information
type ProductInfo struct {
	Name  string
	Price uint32
}

// Describe product
type Item struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

// Describe user order id
type OrderID struct {
	ID int64
}

// Describe cart item
type CartItem struct {
	SKU   uint32
	Count uint16
}

// Describe user cart
type Cart struct {
	OrderID int64
	Items   []CartItem
}

// Describe methods to check the availability of goods in stock
type LomsChecker interface {
	GetStocksBySKU(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, userGoods []CartItem) (OrderID, error)
}

// Describes the method of retrieving product information
type ProductChecker interface {
	GetProductBySKU(ctx context.Context, sku uint32) (ProductInfo, error)
}

// Provide access to the business logic of the service
type Service struct {
	lomsChecker    LomsChecker
	productChecker ProductChecker
}

// Create a new Service instance
func New(lomsChecker LomsChecker, productChecker ProductChecker) *Service {
	return &Service{
		lomsChecker:    lomsChecker,
		productChecker: productChecker,
	}
}
