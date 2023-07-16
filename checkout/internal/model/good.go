// Good models
package model

// Describe good sku
type SKU uint32

// Describe product
type Good struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}
