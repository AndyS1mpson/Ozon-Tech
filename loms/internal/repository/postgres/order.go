// OrderRepository
package postgres

import (
	"context"
	"route256/loms/internal/converter/repository"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/tracer"
	schema "route256/loms/internal/repository/scheme"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type OrderStatus string

const (
	tableNameOrder     = "user_order"
	tableNameOrderItem = "order_item"

	createdStatus  OrderStatus = "new"
	paidStatus     OrderStatus = "payed"
	canceledStatus OrderStatus = "cancelled"
	failedStatus   OrderStatus = "failed"
	waitStatus     OrderStatus = "awaiting payment"
)

// Repository for woring with orders
type OrderRepository struct {
	db *pgxpool.Pool
}

// Create new Order repository instance
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Get order by id
func (r *OrderRepository) GetOrder(ctx context.Context, id model.OrderID) (*model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/get_order")
	defer span.Finish()

	query, agrs, err := psql.Select("user_id").From(tableNameOrder).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for get order"))
	}

	var userID model.UserID
	err = r.db.QueryRow(ctx, query, agrs...).Scan(&userID)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "get order"))
	}

	itemsQuery, agrs, err := psql.Select("sku", "count").From(tableNameOrderItem).Where(sq.Eq{"order_id": id}).ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for get order items"))
	}

	var items []schema.Item
	err = pgxscan.Select(ctx, r.db, &items, itemsQuery, agrs...)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "get order items"))
	}

	orderItems := repository.ToOrderItems(items)

	return &model.Order{
		User:  int64(userID),
		Items: orderItems,
	}, nil
}

// Create user order
func (r *OrderRepository) CreateOrder(ctx context.Context, order model.Order) (model.OrderID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/create_order")
	defer span.Finish()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "begin transaction"))
	}

	// Create the order
	orderID, err := r.createOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback(ctx)
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "create order"))
	}

	err = r.insertOrderItems(ctx, tx, orderID, order.Items)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "insert order items"))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "commit transaction"))
	}

	return orderID, nil
}

// Get order info
func (r *OrderRepository) ListOrder(ctx context.Context, orderID model.OrderID) (model.OrderWithStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/list_order")
	defer span.Finish()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return model.OrderWithStatus{}, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "begin transaction"))
	}

	// Get user id
	user, status, err := r.getOrderUserWithStatus(ctx, tx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "get order"))
	}

	items, err := r.getOrderItems(ctx, tx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "get order items"))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.OrderWithStatus{}, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "commit transaction"))
	}
	return model.OrderWithStatus{
		Status: string(status),
		User:   int64(user),
		Items:  items,
	}, nil
}

// Change order status to paid
func (r *OrderRepository) PayOrder(ctx context.Context, orderID model.OrderID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/pay_order")
	defer span.Finish()

	query := psql.Update(tableNameOrder).Set("status", paidStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for update order"))
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "exec update order"))
	}

	return nil
}

// Change order status to canceled
func (r *OrderRepository) CancelOrder(ctx context.Context, orderID model.OrderID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/cancel_order")
	defer span.Finish()

	query := psql.Update(tableNameOrder).Set("status", canceledStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for update order"))
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "exec update order"))
	}

	return nil
}

// Change status to awaiting payment
func (r *OrderRepository) AwaitPaymentOrder(ctx context.Context, orderID model.OrderID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/await_payment_order")
	defer span.Finish()

	query := psql.Update(tableNameOrder).Set("status", waitStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for update order"))
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "exec update order"))
	}

	return nil
}

// Change status to failed
func (r *OrderRepository) FailOrder(ctx context.Context, orderID model.OrderID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/order/fail_order")
	defer span.Finish()

	query := psql.Update(tableNameOrder).Set("status", failedStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "build query for update order"))
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "exec update order"))
	}

	return nil
}

// Create order object
func (r *OrderRepository) createOrder(ctx context.Context, tx pgx.Tx, order model.Order) (model.OrderID, error) {
	query, args, err := psql.
		Insert(tableNameOrder).
		Columns("user_id", "status").
		Values(order.User, createdStatus).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, errors.Wrap(err, "build query for create order")
	}

	var res model.OrderID
	err = tx.QueryRow(ctx, query, args...).Scan(&res)
	if err != nil {
		return 0, errors.Wrap(err, "exec create order")
	}
	return res, nil
}

// Insert item to order
func (r *OrderRepository) insertOrderItems(ctx context.Context, tx pgx.Tx, orderID model.OrderID, items []model.OrderItem) error {
	query := psql.Insert(tableNameOrderItem).Columns("order_id", "sku", "count")

	for _, item := range items {
		query = query.Values(int64(orderID), item.SKU, item.Count)
	}

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query for insert order items")
	}

	_, err = tx.Exec(ctx, rawSQL, args...)
	if err != nil {
		return errors.Wrap(err, "exec insert order items")
	}

	return nil
}

// Get order's user id with order status
func (r *OrderRepository) getOrderUserWithStatus(ctx context.Context, tx pgx.Tx, orderID model.OrderID) (model.UserID, OrderStatus, error) {
	// Get user id
	query := psql.Select("user_id", "status").From(tableNameOrder).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, "", errors.Wrap(err, "build query for list order")
	}

	var user model.UserID
	var status string
	err = tx.QueryRow(ctx, rawSQL, args...).Scan(&user, &status)
	if err != nil {
		return 0, "", errors.Wrap(err, "get user id")
	}

	return user, OrderStatus(status), nil
}

// Get order items by order id
func (r *OrderRepository) getOrderItems(ctx context.Context, tx pgx.Tx, orderID model.OrderID) ([]model.OrderItem, error) {
	// Get order items
	itemsQuery := psql.Select("sku", "count").From(tableNameOrderItem).Where(sq.Eq{"order_id": orderID})

	rawSQL, args, err := itemsQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query for list order items")
	}

	var items []model.OrderItem
	rows, err := tx.Query(ctx, rawSQL, args...)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "get order items")
	}

	for rows.Next() {
		var item model.OrderItem

		err := rows.Scan(&item.SKU, &item.Count)
		if err != nil {
			return nil, errors.Wrap(err, "scan order items")
		}

		items = append(items, item)
	}

	return items, nil
}
