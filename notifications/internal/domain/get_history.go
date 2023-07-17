package domain

import (
	"context"
	"route256/notifications/internal/model"
	"time"
)

// Get user notifications history
func (s *Service) GetHistoryWithPeriod(ctx context.Context, userID model.UserID, from time.Time, to time.Time) ([]model.OrderStatusMessage, error) {
	cacheKey := CacheKey{
		UserID: int64(userID),
		From:   from,
		To:     to,
	}
	// try to get from cache
	cacheVal, ok := s.cache.Get(cacheKey)
	if ok {
		return cacheVal.items, nil
	}
	// if not in cache, get from db
	messages, err := s.message.GetHistoryWithPeriod(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	if len(messages) != 0 {
		// Save to cache
		s.cache.Set(cacheKey, CacheVal{items: messages})
	}
	return messages, nil
}
