// OrderRepository
package postgres

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/converter/repository"
	"route256/loms/internal/model"
	schema "route256/loms/internal/repository/scheme"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	query, agrs, err := psql.Select("user_id").From(tableNameOrder).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for get order: %s", err)
	}

	var userID model.UserID
	err = r.db.QueryRow(ctx, query, agrs...).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("get order: %s", err)
	}

	itemsQuery, agrs, err := psql.Select("sku", "count").From(tableNameOrderItem).Where(sq.Eq{"order_id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for get order items: %s", err)
	}

	var items []schema.Item
	err = pgxscan.Select(ctx, r.db, &items, itemsQuery, agrs...)
	if err != nil {
		return nil, fmt.Errorf("get order items: %s", err)
	}

	orderItems := repository.ToOrderItems(items)

	return &model.Order{
		User:  int64(userID),
		Items: orderItems,
	}, nil
}

// Create user order
func (r *OrderRepository) CreateOrder(ctx context.Context, order model.Order) (model.OrderID, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %s", err)
	}

	// Create the order
	orderID, err := r.createOrder(ctx, tx, order)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("create order: %s", err)
	}

	err = r.insertOrderItems(ctx, tx, orderID, order.Items)
	if err != nil {
		return 0, fmt.Errorf("insert order items: %s", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %s", err)
	}

	return orderID, nil
}

// Get order info
func (r *OrderRepository) ListOrder(ctx context.Context, orderID model.OrderID) (model.OrderWithStatus, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return model.OrderWithStatus{}, fmt.Errorf("begin transaction: %s", err)
	}

	// Get user id
	user, status, err := r.getOrderUserWithStatus(ctx, tx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, fmt.Errorf("get order: %s", err)
	}

	items, err := r.getOrderItems(ctx, tx, orderID)
	if err != nil {
		return model.OrderWithStatus{}, fmt.Errorf("get order items: %s", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.OrderWithStatus{}, fmt.Errorf("commit transaction: %s", err)
	}
	return model.OrderWithStatus{
		Status: string(status),
		User:   int64(user),
		Items:  items,
	}, nil
}

// Change order status to paid
func (r *OrderRepository) PayOrder(ctx context.Context, orderID model.OrderID) error {
	query := psql.Update(tableNameOrder).Set("status", paidStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update order: %s", err)
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update order: %s", err)
	}

	return nil
}

// Change order status to canceled
func (r *OrderRepository) CancelOrder(ctx context.Context, orderID model.OrderID) error {

	query := psql.Update(tableNameOrder).Set("status", canceledStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update order: %s", err)
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update order: %s", err)
	}

	return nil
}

// Change status to awaiting payment
func (r *OrderRepository) AwaitPaymentOrder(ctx context.Context, orderID model.OrderID) error {
	query := psql.Update(tableNameOrder).Set("status", waitStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update order: %s", err)
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update order: %s", err)
	}

	return nil
}

// Change status to failed
func (r *OrderRepository) FailOrder(ctx context.Context, orderID model.OrderID) error {
	query := psql.Update(tableNameOrder).Set("status", failedStatus).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update order: %s", err)
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update order: %s", err)
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
		return 0, fmt.Errorf("build query for create order: %s", err)
	}

	var res model.OrderID
	err = tx.QueryRow(ctx, query, args...).Scan(&res)
	if err != nil {
		return 0, fmt.Errorf("exec create order: %s", err)
	}
	log.Println("NeW ORDER: ", res)
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
		return fmt.Errorf("build query for insert order items: %s", err)
	}

	_, err = tx.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec insert order items: %s", err)
	}

	return nil
}

// Get order's user id with order status
func (r *OrderRepository) getOrderUserWithStatus(ctx context.Context, tx pgx.Tx, orderID model.OrderID) (model.UserID, OrderStatus, error) {
	// Get user id
	query := psql.Select("user_id", "status").From(tableNameOrder).Where(sq.Eq{"id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, "", fmt.Errorf("build query for list order: %s", err)
	}

	var user model.UserID
	var status string
	err = tx.QueryRow(ctx, rawSQL, args...).Scan(&user, &status)
	if err != nil {
		return 0, "", fmt.Errorf("get user id: %s", err)
	}

	return user, OrderStatus(status), nil
}

// Get order items by order id
func (r *OrderRepository) getOrderItems(ctx context.Context, tx pgx.Tx, orderID model.OrderID) ([]model.OrderItem, error) {
	// Get order items
	itemsQuery := psql.Select("sku", "count").From(tableNameOrderItem).Where(sq.Eq{"order_id": orderID})

	rawSQL, args, err := itemsQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for list order items: %s", err)
	}

	var items []model.OrderItem
	rows, err := tx.Query(ctx, rawSQL, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("get order items: %s", err)
	}

	for rows.Next() {
		var item model.OrderItem

		err := rows.Scan(&item.SKU, &item.Count)
		if err != nil {
			return nil, fmt.Errorf("scan order items: %s", err)
		}

		items = append(items, item)
	}

	return items, nil
}
