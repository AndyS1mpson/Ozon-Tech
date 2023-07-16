//go:build integration

package integrationtest

import (
	"context"
	"route256/checkout/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

const (
	tableNameCart     = "cart"
	tableNameCartItem = "cart_item"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Test getting cart id by user id
func (s *Suite) Test_GetCartByUserID() {
	// Arrange
	userID := 20
	query := psql.Insert(tableNameCart).Columns("user_id").Values(userID).Suffix("RETURNING id")

	rawSQL, args, err := query.ToSql()
	s.Require().NoError(err)

	var result model.UserCartID
	err = s.pg.QueryRow(context.Background(), rawSQL, args...).Scan(&result)
	s.Require().NoError(err)

	// Act
	cartID, err := s.cart.GetCartByUserID(context.Background(), model.UserID(userID))

	// Assert
	s.Require().NoError(err)
	s.Require().NotEmpty(cartID)
}

// Test creating user cart
func (s *Suite) Test_CreateCart() {
	// Arrange
	userID := model.UserID(15)

	// Act
	cartID, err := s.cart.CreateCart(context.Background(), userID)

	// Assert
	s.Require().NoError(err)

	var newCartID int64
	query, args, err := psql.Select("id").From(tableNameCart).Where(sq.Eq{"user_id": userID}).ToSql()
	s.Require().NoError(err)

	err = pgxscan.Get(context.Background(), s.pg, &newCartID, query, args...)
	s.Require().NoError(err)
	s.Require().Equal(newCartID, int64(cartID))
}

// Test updating or adding new item to user cart
func (s *Suite) Test_UpdateOrAddToCart() {
	// Arrange
	userID := model.UserID(15)
	sku := model.SKU(751)
	count := uint16(5)

	insertQuery, args, err := psql.Insert(tableNameCart).Columns("user_id").Values(userID).Suffix("RETURNING id").ToSql()
	s.Require().NoError(err)

	var cartID model.UserCartID
	err = s.pg.QueryRow(context.Background(), insertQuery, args...).Scan(&cartID)
	s.Require().NoError(err)

	// Act
	err = s.cart.UpdateOrAddToCart(context.Background(), cartID, sku, count)

	// Assert
	s.Require().NoError(err)

	var newSKU int64
	selectQuery, args, err := psql.Select("sku").From(tableNameCartItem).Where(sq.Eq{"cart_id": cartID}).ToSql()
	s.Require().NoError(err)

	err = pgxscan.Get(context.Background(), s.pg, &newSKU, selectQuery, args...)
	s.Require().NoError(err)
	s.Require().Equal(newSKU, int64(sku))
}

// Test removing item from user cart
func (s *Suite) Test_DeleteFromCart() {
	// Arrange
	userID := model.UserID(15)
	sku := model.SKU(751)
	count := uint16(5)
	removedCount := uint16(2)

	insertQuery, args, err := psql.Insert(tableNameCart).Columns("user_id").Values(userID).Suffix("RETURNING id").ToSql()
	s.Require().NoError(err)

	var cartID model.UserCartID
	err = s.pg.QueryRow(context.Background(), insertQuery, args...).Scan(&cartID)
	s.Require().NoError(err)

	insertItemQuery, args, err := psql.Insert(tableNameCartItem).Columns("cart_id", "sku", "count").Values(cartID, sku, count).ToSql()
	s.Require().NoError(err)

	_, err = s.pg.Exec(context.Background(), insertItemQuery, args...)
	s.Require().NoError(err)

	// Act
	err = s.cart.DeleteFromCart(context.Background(), cartID, sku, removedCount)

	// Assert
	s.Require().NoError(err)

	var stock int
	selectQuery, args, err := psql.Select("count").From(tableNameCartItem).Where(sq.Eq{"cart_id": cartID, "sku": sku}).ToSql()
	s.Require().NoError(err)

	err = pgxscan.Get(context.Background(), s.pg, &stock, selectQuery, args...)
	s.Require().NoError(err)
	s.Require().Equal(stock, int(count-removedCount))
}

// Test getting all items from user cart
func (s *Suite) Test_ListCart() {
	// Arrange
	userID := model.UserID(15)
	sku1 := model.SKU(751)
	count1 := uint16(5)
	sku2 := model.SKU(760)
	count2 := uint16(5)

	insertQuery, args, err := psql.Insert(tableNameCart).Columns("user_id").Values(userID).Suffix("RETURNING id").ToSql()
	s.Require().NoError(err)

	var cartID model.UserCartID
	err = s.pg.QueryRow(context.Background(), insertQuery, args...).Scan(&cartID)
	s.Require().NoError(err)

	insertItemQuery, args, err := psql.
		Insert(tableNameCartItem).
		Columns("cart_id", "sku", "count").
		Values(cartID, sku1, count1).
		Values(cartID, sku2, count2).
		ToSql()
	s.Require().NoError(err)

	_, err = s.pg.Exec(context.Background(), insertItemQuery, args...)
	s.Require().NoError(err)

	// Act
	cart, err := s.cart.ListCart(context.Background(), cartID)

	// Assert
	s.Require().NoError(err)
	s.Require().Len(cart, 2)
	s.Require().ElementsMatch(cart, []model.CartItem{
		{SKU: uint32(sku1), Count: count1},
		{SKU: uint32(sku2), Count: count1},
	})
}
