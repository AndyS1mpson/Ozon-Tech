// CreateOrder
package notifications

import (
	"context"
	"route256/notifications/internal/converter/server"
	"route256/notifications/internal/model"
	"route256/notifications/pkg/notifications_v1"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateOrder controller
func (s *Server) GetHistoryWithPeriod(ctx context.Context, req *notifications_v1.GetHistoryWithPeriodRequest) (*notifications_v1.GetHistoryWithPeriodResponse, error) {
	messages, err := s.service.GetHistoryWithPeriod(
		ctx,
		model.UserID(req.GetUser()),
		time.Date(int(req.From.GetYear()), time.Month(req.From.GetMonth()), int(req.From.GetDay()), 0, 0, 0, 0, time.UTC),
		time.Date(int(req.To.GetYear()), time.Month(req.To.GetMonth()), int(req.To.GetDay()), 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		return &notifications_v1.GetHistoryWithPeriodResponse{}, status.Errorf(codes.Internal, err.Error())
	}
	return server.ListMessagesToResp(messages), nil
}
