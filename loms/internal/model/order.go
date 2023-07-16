// Define order's DTO for domain layer
package model

// Definte user identifier
type UserID int64

// Define order identifier
type OrderID int64

// Define user order good
type OrderItem struct {
	SKU   uint32
	Count uint16
}

// Define user order
type Order struct {
	User  int64
	Items []OrderItem
}

type OrderStatus string

const (
	CreatedStatus  OrderStatus = "new"
	PaidStatus     OrderStatus = "payed"
	CanceledStatus OrderStatus = "cancelled"
	FailedStatus   OrderStatus = "failed"
	WaitStatus     OrderStatus = "awaiting payment"
)

// Define user order info with status
type OrderWithStatus struct {
	Status string
	User   int64
	Items  []OrderItem
}
