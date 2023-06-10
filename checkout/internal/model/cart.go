// User cart models
package model

// Describe user id
type UserID int64

// Describe order id
type OrderID int64

// Describe user cart id
type UserCartID int64

// Describe user cart
type Cart struct {
	ID     int64
	UserID int64
}

// Describe cart item
type CartItem struct {
	SKU   uint32
	Count uint16
}

// Describe user cart with total price
type UserCartWithTotal struct {
	Items      []Good
	TotalPrice uint32
}
