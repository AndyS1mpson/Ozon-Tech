// Cart item table definition
package schema

// Describe cart item table in postgres
type CartItem struct {
	CartID int64  `db:"cart_id"`
	SKU    uint32 `db:"sku"`
	Count  uint16 `db:"count"`
}
