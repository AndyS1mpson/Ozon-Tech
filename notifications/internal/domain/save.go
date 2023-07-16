package domain

import (
	"context"
	"route256/notifications/internal/model"
)

// Save message to database
func (s *Service) Save(ctx context.Context, message model.OrderStatusMessage) (model.MessageID, error) {
	id, err := s.message.Save(ctx, message)
	if err != nil {
		return 0, err
	}
	return id, nil
}
