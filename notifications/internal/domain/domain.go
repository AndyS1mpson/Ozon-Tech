package domain

import (
	"context"
	"route256/notifications/internal/model"
	"time"
)

// Describe repository for working with messages
type MessageRepository interface {
	Save(ctx context.Context, message model.OrderStatusMessage) (model.MessageID, error)
	GetHistoryWithPeriod(ctx context.Context, userID model.UserID, from time.Time, to time.Time) ([]model.OrderStatusMessage, error)
}

// Describe a client that will notify user about events
type Notifier interface {
	SendMessage(message string) error
}

// Describe a cache key
type CacheKey struct {
	UserID int64
	From   time.Time
	To     time.Time
}

// Describe a cache value
type CacheVal struct {
	items []model.OrderStatusMessage
}

// Describe cache
type Cacher interface {
	Get(key CacheKey) (CacheVal, bool)
	Set(key CacheKey, value CacheVal) error
	Delete(key CacheKey)
	Clear()
	Count() int
}

// Implement business-logic
type Service struct {
	message  MessageRepository
	notifier Notifier
	cache    Cacher
}

// Create new service instance
func NewService(message MessageRepository, notifier Notifier, cache Cacher) *Service {
	return &Service{
		message:  message,
		notifier: notifier,
		cache:    cache,
	}
}
