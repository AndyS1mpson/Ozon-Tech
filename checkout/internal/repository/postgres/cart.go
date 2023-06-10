package postgres

import (
	"context"
	"fmt"
	"log"
	"route256/checkout/internal/converter/repository"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const (
	tableNameCart     = "cart"
	tableNameCartItem = "cart_item"
)

// Define Cart repository
type CartRepository struct {
	db *pgxpool.Pool
}

// Create a new Cart repository instance
func New(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// Get user cart id
func (r *CartRepository) GetCartByUserID(ctx context.Context, userID model.UserID) (model.UserCartID, error) {
	query := psql.Select("*").From(tableNameCart).Where(sq.Eq{"user_id": userID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %s", err)
	}

	var result schema.Cart

	err = pgxscan.Get(ctx, r.db, &result, rawSQL, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query for filter: %s", err)
	}

	return model.UserCartID(result.ID), nil
}

// Create cart for user
func (r *CartRepository) CreateCart(ctx context.Context, user model.UserID) (model.UserCartID, error) {
	query := psql.Insert(tableNameCart).Columns("user_id").Values(int64(user)).Suffix("RETURNING id")

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for create user cart: %s", err)
	}

	var result model.UserCartID
	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(&result)
	log.Println("NEW ITEM: ", result)
	if err != nil {
		return 0, fmt.Errorf("exec insert item: %s", err)
	}

	return result, nil
}

// Add item to user cart or update if it already exists
func (r *CartRepository) UpdateOrAddToCart(ctx context.Context, cart int64, sku uint32, count uint16) error {
	query := psql.Select("cart_id").From(tableNameCartItem).Where(sq.Eq{"cart_id": cart, "sku": sku})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %s", err)
	}

	var cartID int64

	err = r.db.QueryRow(ctx, rawSQL, args...).Scan(&cartID)

	if err == pgx.ErrNoRows {
		// If item not exists in cart
		insertQuery, args, err := psql.
			Insert(tableNameCartItem).
			Columns("cart_id", "sku", "count").
			Values(cart, sku, count).
			ToSql()

		if err != nil {
			return fmt.Errorf("build insert query: %s", err)
		}

		_, err = r.db.Exec(ctx, insertQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	} else if err == nil {
		// If item already exists in cart
		updateQuery, args, err := psql.
			Update(tableNameCartItem).
			Set("count", sq.Expr("count + ?", count)).
			Where(sq.Eq{"cart_id": cart, "sku": sku}).
			ToSql()
		if err != nil {
			return fmt.Errorf("failed to build update query: %s", err)
		}

		_, err = r.db.Exec(ctx, updateQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check item existence: %w", err)
	}

	return nil
}

// Remove item from user cart
func (r *CartRepository) DeleteFromCart(ctx context.Context, cart int64, sku uint32, count uint16) error {
	query := psql.
		Update(tableNameCartItem).
		Set("count", sq.Expr("count - ?", count)).
		Where(sq.Eq{"cart_id": cart})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %s", err)
	}

	_, err = r.db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// Get list of items from user cart
func (r *CartRepository) ListCart(ctx context.Context, cart model.UserCartID) ([]model.CartItem, error) {
	query := psql.Select("*").From(tableNameCartItem).Where(sq.Eq{"cart_id": cart})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %s", err)
	}

	var result []schema.CartItem

	err = pgxscan.Select(ctx, r.db, &result, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query for filter: %s", err)
	}

	cartItems := repository.ToCartItems(result)
	return cartItems, nil
}
