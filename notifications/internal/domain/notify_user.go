package domain

import (
	"context"
	"fmt"
	"route256/notifications/internal/model"
)

// Notify user about order status
func (s *Service) NotifyUser(ctx context.Context, message model.OrderStatusMessage) error {
	strMessage := fmt.Sprintf(
		"Your id: %v, your order: %v have the following status: %v. More information: %v",
		message.UserID,
		message.OrderID,
		message.Status,
		message.Message,
	)
	err := s.notifier.SendMessage(strMessage)
	if err != nil {
		return err
	}
	return nil
}
