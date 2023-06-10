// gRPC server
package cart

import (
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/cart_v1"
)

// Define gRPC server
type Server struct {
	cart_v1.UnimplementedCartServer
	service *domain.Service
}

// Create new server instance
func New(service *domain.Service) *Server {
	return &Server{service: service}
}
