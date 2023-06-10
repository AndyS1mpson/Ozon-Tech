// gRPC server
package loms

import (
	"route256/loms/internal/domain"
	"route256/loms/pkg/loms_v1"
)

// Define gRPC server
type Server struct {
	loms_v1.UnimplementedLomsServer
	service *domain.Service
}

// Create new Server instance
func New(service *domain.Service) *Server {
	return &Server{service: service}
}
