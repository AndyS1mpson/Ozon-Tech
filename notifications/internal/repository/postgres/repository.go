package postgres

import (
	"context"
	"route256/notifications/internal/model"
	"route256/notifications/internal/pkg/tracer"
	"route256/notifications/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const (
	tableNameMessage = "message"
)

// Define user status history repository
type MessageRepository struct {
	db *pgxpool.Pool
}

// Create a new MessageRepository instance
func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

// Save message
func (m *MessageRepository) Save(ctx context.Context, message model.OrderStatusMessage) (model.MessageID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/cart/update_or_add_item")
	defer span.Finish()

	// If item not exists in cart
	insertQuery, args, err := psql.
		Insert(tableNameMessage).
		Columns("user_id", "order_id", "status", "message").
		Values(message.UserID, message.OrderID, message.Status, message.Message).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build insert query"))
	}

	var messageID model.MessageID

	err = m.db.QueryRow(ctx, insertQuery, args...).Scan(&messageID)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to insert item"))
	}

	return messageID, nil
}

// Get user messages from period
func (m *MessageRepository) GetHistoryWithPeriod(ctx context.Context, userID model.UserID, from time.Time, to time.Time) ([]model.OrderStatusMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/cart/update_or_add_item")
	defer span.Finish()

	query := psql.Select("*").From(tableNameMessage).Where(sq.Eq{"user_id": userID}).Where(sq.GtOrEq{"created_date": from}).Where(sq.LtOrEq{"created_date": to})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query"))
	}

	var messages []schema.MessageItem

	err = pgxscan.Select(ctx, m.db, &messages, rawSQL, args...)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "exec query for filter"))
	}

	var result []model.OrderStatusMessage

	for _, message := range messages {
		result = append(result, model.OrderStatusMessage{
			UserID:  model.UserID(message.UserID),
			OrderID: model.OrderID(message.OrderID),
			Status:  model.OrderStatus(message.Status),
			Message: message.Message,
		})
	}

	return result, nil
}
