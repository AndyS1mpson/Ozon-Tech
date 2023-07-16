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

func Test_ListCart(t *testing.T) {
	t.Parallel()

	t.Run("success, items: [3]", func(t *testing.T) {
		t.Parallel()
		// Arrange
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

		// fill product service data
		goods := make([]model.Good, userItemsCount)
		totalPrice := uint32(0)
		for i := 0; i < userItemsCount; i++ {
			require.NoError(t, gofakeit.Struct(&goods[i]))
			goods[i].SKU = items[i].SKU
			goods[i].Count = items[i].Count

			totalPrice += uint32(goods[i].Count) * goods[i].Price
		}
		product.On("GetProducts", mock.Anything, items).Return(goods, nil).Once()

		service := New(loms, product, cartRepository)

		// Act
		userCart, err := service.ListCart(context.Background(), userID)

		// Assert
		require.NoError(t, err)
		require.Equal(t, goods, userCart.Items, "cart items")
		require.Equal(t, userCart.TotalPrice, totalPrice, "cart total price")
		require.ElementsMatch(t, goods, userCart.Items)
	})

	t.Run("success if the user's shopping cart did not exist", func(t *testing.T) {
		t.Parallel()
		// Arrange
		errStub := errors.New("stub")
		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		userID := model.UserID(1)
		userCartID := model.UserCartID(1)

		// fill repository mock data
		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(model.UserCartID(0), errStub).Once()
		cartRepository.On("CreateCart", mock.Anything, userID).Return(userCartID, nil).Once()

		// fill product service data
		service := New(loms, product, cartRepository)

		// Act
		userCart, err := service.ListCart(context.Background(), userID)

		// Assert
		require.NoError(t, err)
		require.Equal(t, model.UserCartWithTotal{}, userCart, "empty user cart")
	})

	t.Run("error while creating cart from storage", func(t *testing.T) {
		t.Parallel()
		// Arrange
		errStub := errors.New("stub")
		userID := model.UserID(1)

		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)
		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(model.UserCartID(0), errStub).Once()
		cartRepository.On("CreateCart", mock.Anything, userID).Return(model.UserCartID(0), errStub).Once()

		service := New(loms, product, cartRepository)
		// Act
		_, err := service.ListCart(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error while cart items from storage", func(t *testing.T) {
		t.Parallel()
		// Arrange
		errStub := errors.New("stub")
		userID := model.UserID(1)
		cartID := model.UserCartID(1)

		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)
		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(cartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, cartID).Return(nil, errStub).Once()

		service := New(loms, product, cartRepository)
		// Act
		_, err := service.ListCart(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error can not get products", func(t *testing.T) {
		t.Parallel()
		// Arrange
		userItemsCount := 3
		errStub := errors.New("stub")
		userID := model.UserID(1)
		cartID := model.UserCartID(1)

		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		// fill repository mock data
		items := make([]model.CartItem, userItemsCount)
		for i := 0; i < userItemsCount; i++ {
			require.NoError(t, gofakeit.Struct(&items[i]))
		}

		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(cartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, cartID).Return(items, nil).Once()
		product.On("GetProducts", mock.Anything, items).Return(nil, errStub).Once()

		service := New(loms, product, cartRepository)

		// Act
		_, err := service.ListCart(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errStub)
	})

	t.Run("error can not get all products", func(t *testing.T) {
		t.Parallel()
		// Arrange
		userItemsCount := 3
		userID := model.UserID(1)
		cartID := model.UserCartID(1)

		loms := mocks.NewLomsChecker(t)
		product := mocks.NewProductChecker(t)
		cartRepository := mocks.NewCartRepository(t)

		// fill repository mock data
		items := make([]model.CartItem, userItemsCount)
		for i := 0; i < userItemsCount; i++ {
			require.NoError(t, gofakeit.Struct(&items[i]))
		}

		cartRepository.On("GetCartByUserID", mock.Anything, userID).Return(cartID, nil).Once()
		cartRepository.On("ListCart", mock.Anything, cartID).Return(items, nil).Once()

		// fill product service data
		goods := make([]model.Good, userItemsCount-1)
		totalPrice := uint32(0)
		for i := 0; i < userItemsCount-1; i++ {
			require.NoError(t, gofakeit.Struct(&goods[i]))
			goods[i].SKU = items[i].SKU
			goods[i].Count = items[i].Count

			totalPrice += uint32(goods[i].Count) * goods[i].Price
		}
		product.On("GetProducts", mock.Anything, items).Return(goods, nil).Once()

		service := New(loms, product, cartRepository)

		// Act
		_, err := service.ListCart(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, ErrGetProductsInfo)
	})
}
