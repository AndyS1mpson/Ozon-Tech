// Message history table definition
package schema

import "time"

// Describe message item table in postgres
type MessageItem struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	OrderID     int64     `db:"order_id"`
	Status      string    `db:"status"`
	Message     string    `db:"message"`
	CreatedDate time.Time `db:"created_date"`
}
