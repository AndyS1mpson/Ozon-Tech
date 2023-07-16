package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/domain/mocks"
	"route256/checkout/internal/model"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Purchase(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userItemsCount := 3
		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		userID := model.UserID(1)
		userCartID := model.UserCartID(1)
		orderID := model.OrderID(1)

		// fill repository mock data
		items := make([]model.CartItem, userItemsCount)
		for i := 0; i < userItemsCount; i++ {
			require.NoError(t, gofakeit.Struct(&items[i]))
		}
		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(userCartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, userCartID).Return(items, nil).Once()
		loms.On("CreateOrder", mock.Anything, userID, items).Return(orderID, nil).Once()

		service := New(loms, product, cartRepository)

		// Act
		id, err := service.Purchase(context.Background(), userID)
		// Assert
		require.NoError(t, err)
		require.Equal(t, id, orderID)
	})

	t.Run("error user have empty cart", func(t *testing.T) {
		t.Parallel()

		// Arrange
		errStub := errors.New("stub")
		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		userID := model.UserID(1)

		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(model.UserCartID(0), errStub).Once()

		service := New(loms, product, cartRepository)

		// Act
		_, err := service.Purchase(context.Background(), userID)
		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error get user cart items", func(t *testing.T) {
		t.Parallel()

		// Arrange
		errStub := errors.New("stub")
		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		userID := model.UserID(1)
		userCartID := model.UserCartID(1)

		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(userCartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, userCartID).Return(nil, errStub).Once()

		service := New(loms, product, cartRepository)

		// Act
		_, err := service.Purchase(context.Background(), userID)
		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("purchase order", func(t *testing.T) {
		t.Parallel()

		// Arrange
		errStub := errors.New("stub")
		userItemsCount := 3
		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		userID := model.UserID(1)
		userCartID := model.UserCartID(1)

		// fill repository mock data
		items := make([]model.CartItem, userItemsCount)
		for i := 0; i < userItemsCount; i++ {
			require.NoError(t, gofakeit.Struct(&items[i]))
		}
		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(userCartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, userCartID).Return(items, nil).Once()
		loms.On("CreateOrder", mock.Anything, userID, items).Return(model.OrderID(0), errStub).Once()

		service := New(loms, product, cartRepository)

		// Act
		_, err := service.Purchase(context.Background(), userID)
		// Assert
		require.ErrorIs(t, err, errStub)
	})
}
