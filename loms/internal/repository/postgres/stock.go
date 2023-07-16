// StockRepository
package postgres

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/tracer"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	tableNameStock         = "stock"
	tableNameReservedStock = "reservation_stock"
)

// Repository for working with stock and reserved items
type StockRepository struct {
	db *pgxpool.Pool
}

// Create new stock repository instance
func NewStockRepository(db *pgxpool.Pool) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

// Get stocks where there is a free product
func (s *StockRepository) GetAvailableStocks(ctx context.Context, sku model.SKU) ([]model.Stock, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/stocks/get_available_stocks")
	defer span.Finish()

	query, args, err := psql.
		Select("warehouse_id", "count").
		From(tableNameStock).
		Where(sq.Eq{"sku": uint32(sku)}).
		ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to build query"))
	}

	rows, err := s.db.Query(ctx, query, args...)
	defer rows.Close()

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to get stocks"))
	}
	stocks := make([]model.Stock, 0)

	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.WarehouseID, &stock.Count)
		if err != nil {
			return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to scan stock"))
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// Reserve Item
func (r *StockRepository) Reserve(ctx context.Context, orderID model.OrderID, sku model.SKU, stock model.Stock) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/stocks/reserve")
	defer span.Finish()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to begin transaction"))
	}

	err = r.addItemToReserve(ctx, tx, orderID, sku, stock)
	if err != nil {
		tx.Rollback(ctx)
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to add item to reserve"))
	}

	err = r.removeItemFromStock(ctx, tx, sku, stock)
	if err != nil {
		tx.Rollback(ctx)
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to remove item from stock"))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "commit transaction"))
	}

	return nil
}

// Unreserve Item
func (r *StockRepository) Unreserve(ctx context.Context, orderID model.OrderID, sku model.SKU) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/stocks/unreserve")
	defer span.Finish()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to begin transaction"))
	}

	unresItems, err := r.removeItemFromReserve(ctx, tx, orderID, sku)
	if err != nil {
		tx.Rollback(ctx)
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to remove item from reserve"))
	}

	for _, item := range unresItems {
		err = r.addItemToStock(ctx, tx, sku, item)
		if err != nil {
			tx.Rollback(ctx)
			return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to add item to stock"))
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, errors.Wrap(err, "commit transaction"))
	}

	return nil
}

// Remove items from reserve
func (r *StockRepository) WriteOffOrderItems(ctx context.Context, orderID model.OrderID) ([]model.Stock, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/stocks/write_off_order_items")
	defer span.Finish()

	// Get item warehouse and count for return to stock
	selectQuery, agrs, err := psql.
		Select("warehouse_id", "count").
		From(tableNameReservedStock).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to build query"))
	}

	var resultSQL []struct {
		WarehouseID int64 `db:"warehouse_id"`
		Count       int64 `db:"count"`
	}
	err = pgxscan.Select(ctx, r.db, &resultSQL, selectQuery, agrs...)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to remove item from reserve"))
	}

	result := make([]model.Stock, len(resultSQL))
	for i, v := range resultSQL {
		result[i] = model.Stock{
			WarehouseID: v.WarehouseID,
			Count:       uint64(v.Count),
		}
	}

	// Remove item from reserve
	query, args, err := psql.
		Delete(tableNameReservedStock).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to build query"))
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "failed to update item"))
	}

	return result, nil
}

// Reserve item for stock
func (r *StockRepository) addItemToReserve(ctx context.Context, tx pgx.Tx, orderID model.OrderID, sku model.SKU, stock model.Stock) error {
	insertQuery, args, err := psql.
		Insert(tableNameReservedStock).
		Columns("order_id", "sku", "warehouse_id", "count").
		Values(orderID, sku, stock.WarehouseID, stock.Count).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	_, err = tx.Exec(ctx, insertQuery, args...)
	if err != nil {
		return errors.Wrap(err, "failed to add item to reserve")
	}

	return nil
}

// Remove item from stock
func (r *StockRepository) removeItemFromReserve(ctx context.Context, tx pgx.Tx, orderID model.OrderID, sku model.SKU) ([]model.Stock, error) {
	// Get item warehouse and count for return to stock
	selectQuery, agrs, err := psql.Select("warehouse_id", "count").From(tableNameReservedStock).Where(sq.Eq{"order_id": orderID, "sku": sku}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	var resultSQL []struct {
		WarehouseID int64 `db:"warehouse_id"`
		Count       int64 `db:"count"`
	}
	err = pgxscan.Select(ctx, r.db, &resultSQL, selectQuery, agrs...)

	if err != nil {
		return nil, errors.Wrap(err, "failed to remove item from reserve")
	}
	result := make([]model.Stock, len(resultSQL))
	for i, v := range resultSQL {
		result[i] = model.Stock{
			WarehouseID: v.WarehouseID,
			Count:       uint64(v.Count),
		}
	}

	// Remove item from reserve
	query, args, err := psql.
		Delete(tableNameReservedStock).
		Where(sq.Eq{"sku": sku, "order_id": orderID}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update item")
	}

	return result, nil
}

// Add item to stock
func (r *StockRepository) addItemToStock(ctx context.Context, tx pgx.Tx, sku model.SKU, stock model.Stock) error {
	query, args, err := psql.Select("warehouse_id").
		From(tableNameStock).
		Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}

	var warehouseID model.WarehouseID
	err = tx.QueryRow(ctx, query, args...).Scan(&warehouseID)
	if err == pgx.ErrNoRows {
		insertQuery, args, err := psql.
			Insert(tableNameStock).
			Columns("sku", "warehouse_id", "count").
			Values(sku, stock.WarehouseID, stock.Count).
			ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		_, err = tx.Exec(ctx, insertQuery, args...)
		if err != nil {
			return errors.Wrap(err, "failed to add item to reserve")
		}

	} else if err == nil {
		updateQuery, args, err := psql.
			Update(tableNameStock).
			Set("count", sq.Expr("count + ?", stock.Count)).
			Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).
			ToSql()

		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		_, err = tx.Exec(ctx, updateQuery, args...)
		if err != nil {
			return errors.Wrap(err, "failed to update item")
		}
	} else {
		return errors.Wrap(err, "failed to exec request")
	}

	return nil
}

// Remove amount of item from stock
func (r *StockRepository) removeItemFromStock(ctx context.Context, tx pgx.Tx, sku model.SKU, stock model.Stock) error {
	selectQuery, args, err := psql.Select("count").From(tableNameStock).Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build query")
	}
	var count uint64
	err = tx.QueryRow(ctx, selectQuery, args...).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "not item in warehouse")
	}

	if count-stock.Count == 0 {
		query, args, err := psql.
			Delete(tableNameStock).
			Where(sq.Eq{"sku": uint32(sku), "warehouse_id": stock.WarehouseID}).
			ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return errors.Wrap(err, "failed to update item")
		}
	} else {
		query, args, err := psql.
			Update(tableNameStock).
			Set("count", sq.Expr("count -?", stock.Count)).
			Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).
			ToSql()

		if err != nil {
			return errors.Wrap(err, "failed to build query")
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return errors.Wrap(err, "failed to update item")
		}
	}

	return nil
}
