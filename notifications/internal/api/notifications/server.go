// gRPC server
package notifications

import (
	"route256/notifications/internal/domain"
	"route256/notifications/pkg/notifications_v1"
)

// Define gRPC server
type Server struct {
	notifications_v1.UnimplementedNotificationsServer
	service *domain.Service
}

// Create new Server instance
func New(service *domain.Service) *Server {
	return &Server{service: service}
}
