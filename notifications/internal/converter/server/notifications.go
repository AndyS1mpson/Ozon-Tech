// Converters for the presentation layer
package server

import (
	"route256/notifications/internal/model"
	"route256/notifications/pkg/notifications_v1"
)

// Convert message to response object
func MessageToRes(item model.OrderStatusMessage) *notifications_v1.Message {
	return &notifications_v1.Message{
		UserId:  int64(item.UserID),
		OrderId: int64(item.OrderID),
		Status:  string(item.Status),
		Message: item.Message,
	}
}

// Convert messages info to response object
func ListMessagesToResp(messages []model.OrderStatusMessage) *notifications_v1.GetHistoryWithPeriodResponse {
	items := []*notifications_v1.Message{}
	for _, message := range messages {
		items = append(items, MessageToRes(message))
	}

	return &notifications_v1.GetHistoryWithPeriodResponse{
		Messages: items,
	}
}
