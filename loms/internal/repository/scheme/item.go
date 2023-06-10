// Cart item table definition
package schema

// Describe order item table in postgres
type Item struct {
	orderID int64  `db:"order_id"`
	SKU    uint32 `db:"sku"`
	Count   uint16 `db:"count"`
}
