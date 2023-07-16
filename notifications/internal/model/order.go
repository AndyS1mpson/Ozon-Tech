// Define order's DTO for domain layer
package model

// Define order identifier
type OrderID int64

// Define user identifier
type UserID int64

type OrderStatus string

type MessageID int64

// Define order status message
type OrderStatusMessage struct {
	UserID  UserID
	OrderID OrderID
	Status  OrderStatus
	Message string
}

const (
	CreatedStatus  OrderStatus = "new"
	PaidStatus     OrderStatus = "payed"
	CanceledStatus OrderStatus = "cancelled"
	FailedStatus   OrderStatus = "failed"
	WaitStatus     OrderStatus = "awaiting payment"
)
