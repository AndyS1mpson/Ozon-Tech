// StockRepository
package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	query, args, err := psql.
		Select("warehouse_id", "count").
		From(tableNameStock).
		Where(sq.Eq{"sku": uint32(sku)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}
	stocks := make([]model.Stock, 0)

	for rows.Next() {
		var stock model.Stock
		err := rows.Scan(&stock.WarehouseID, &stock.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// Reserve Item
func (r *StockRepository) Reserve(ctx context.Context, orderID model.OrderID, sku model.SKU, stock model.Stock) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = r.addItemToReserve(ctx, tx, orderID, sku, stock)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to add item to reserve: %w", err)
	}

	err = r.removeItemFromStock(ctx, tx, sku, stock)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to remove item from stock: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %s", err)
	}

	return nil
}

// Unreserve Item
func (r *StockRepository) Unreserve(ctx context.Context, orderID model.OrderID, sku model.SKU) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	unresItems, err := r.removeItemFromReserve(ctx, tx, orderID, sku)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to remove item from reserve: %w", err)
	}

	for _, item := range unresItems {
		err = r.addItemToStock(ctx, tx, sku, item)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to add item to stock: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %s", err)
	}

	return nil
}

// Remove items from reserve
func (r *StockRepository) WriteOffOrderItems(ctx context.Context, orderID model.OrderID) ([]model.Stock, error) {
	// Get item warehouse and count for return to stock
	selectQuery, agrs, err := psql.
		Select("warehouse_id", "count").
		From(tableNameReservedStock).
		Where(sq.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var resultSQL []struct {
		WarehouseID int64 `db:"warehouse_id"`
		Count       int64 `db:"count"`
	}
	err = pgxscan.Select(ctx, r.db, &resultSQL, selectQuery, agrs...)

	if err != nil {
		return nil, fmt.Errorf("failed to remove item from reserve: %w", err)
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
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
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
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.Exec(ctx, insertQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to add item to reserve: %w", err)
	}

	return nil
}

// Remove item from stock
func (r *StockRepository) removeItemFromReserve(ctx context.Context, tx pgx.Tx, orderID model.OrderID, sku model.SKU) ([]model.Stock, error) {
	// Get item warehouse and count for return to stock
	selectQuery, agrs, err := psql.Select("warehouse_id", "count").From(tableNameReservedStock).Where(sq.Eq{"order_id": orderID, "sku": sku}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var resultSQL []struct {
		WarehouseID int64 `db:"warehouse_id"`
		Count       int64 `db:"count"`
	}
	err = pgxscan.Select(ctx, r.db, &resultSQL, selectQuery, agrs...)

	if err != nil {
		return nil, fmt.Errorf("failed to remove item from reserve: %w", err)
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
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
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
		return fmt.Errorf("failed to build query: %w", err)
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
			return fmt.Errorf("failed to build query: %w", err)
		}

		_, err = tx.Exec(ctx, insertQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to add item to reserve: %w", err)
		}

	} else if err == nil {
		updateQuery, args, err := psql.
			Update(tableNameStock).
			Set("count", sq.Expr("count + ?", stock.Count)).
			Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).
			ToSql()

		if err != nil {
			return fmt.Errorf("failed to build query: %w", err)
		}

		_, err = tx.Exec(ctx, updateQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	} else {
		return fmt.Errorf("failed to exec request: %w", err)
	}

	return nil
}

// Remove amount of item from stock
func (r *StockRepository) removeItemFromStock(ctx context.Context, tx pgx.Tx, sku model.SKU, stock model.Stock) error {
	selectQuery, args, err := psql.Select("count").From(tableNameStock).Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	var count uint64
	err = tx.QueryRow(ctx, selectQuery, args...).Scan(&count)
	if err != nil {
		return fmt.Errorf("not item in warehouse: %w", err)
	}

	if count-stock.Count == 0 {
		query, args, err := psql.
			Delete(tableNameStock).
			Where(sq.Eq{"sku": uint32(sku), "warehouse_id": stock.WarehouseID}).
			ToSql()
		if err != nil {
			return fmt.Errorf("failed to build query: %w", err)
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	} else {
		query, args, err := psql.
			Update(tableNameStock).
			Set("count", sq.Expr("count -?", stock.Count)).
			Where(sq.Eq{"sku": sku, "warehouse_id": stock.WarehouseID}).
			ToSql()

		if err != nil {
			return fmt.Errorf("failed to build query: %w", err)
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	}

	return nil
}
