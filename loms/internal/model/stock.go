// Define stocks DTO for domain layer
package model

// Define item identifier
type SKU uint32

// Define warehouse id
type WarehouseID int64

// Describe information about the amount of goods in stock
type Stock struct {
	WarehouseID int64
	Count       uint64
}
