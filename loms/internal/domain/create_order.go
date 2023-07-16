// Creating a custom order
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

// Create user order
func (s *Service) CreateOrder(ctx context.Context, order model.Order) (model.OrderID, error) {

	// Create order
	orderID, err := s.order.CreateOrder(ctx, order)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create order")
	}

	err = s.notifier.SendMessage(OrderStatusNotification{
		UserId: model.UserID(order.User),
		OrderID: orderID,
		Status:  model.CreatedStatus,
		Message: "order is created",
	})
	if err != nil {
		return 0, errors.Wrap(err, "Can not notify about order status")
	}

	// Reserve items from order
	for _, v := range order.Items {
		stocks, err := s.stock.GetAvailableStocks(ctx, model.SKU(v.SKU))
		if err != nil {
			s.order.FailOrder(ctx, orderID)
			err = s.notifier.SendMessage(OrderStatusNotification{
				UserId: model.UserID(order.User),
				OrderID: orderID,
				Status:  model.FailedStatus,
				Message: fmt.Sprintf("no free item: %v for your order", v.SKU),
			})
			if err != nil {
				return 0, errors.Wrap(err, "Can not notify about order status")
			}
			return orderID, errors.Wrap(err, "no available item stocks")
		}

		// Reserve the required quantity from available stocks
		remainingQuantity := v.Count
		for _, stock := range stocks {
			var countToAdd uint64
			if stock.Count >= uint64(remainingQuantity) {
				stock.Count -= uint64(remainingQuantity)
				countToAdd = uint64(remainingQuantity)
				remainingQuantity = 0
			} else {
				remainingQuantity -= uint16(stock.Count)
				countToAdd = stock.Count
				stock.Count = 0
			}

			err = s.stock.Reserve(
				ctx,
				orderID,
				model.SKU(v.SKU),
				model.Stock{
					WarehouseID: stock.WarehouseID,
					Count:       countToAdd,
				},
			)
			if err != nil {
				return orderID, errors.Wrap(err, "failed to update stock count")
			}

			if remainingQuantity == 0 {
				break
			}
		}

		if remainingQuantity != 0 {
			s.order.FailOrder(ctx, orderID)
			err = s.notifier.SendMessage(OrderStatusNotification{
				UserId: model.UserID(order.User),
				OrderID: orderID,
				Status:  model.FailedStatus,
				Message: fmt.Sprintf("no free item: %v for your order", v.SKU),
			})
			if err != nil {
				return 0, errors.Wrap(err, "Can not notify about order status")
			}
			return orderID, errors.Wrap(err, "not enough available stocks")
		}
	}

	err = s.order.AwaitPaymentOrder(ctx, orderID)
	if err != nil {
		return orderID, errors.Wrap(err, "failed to await payment order")
	}

	err = s.notifier.SendMessage(OrderStatusNotification{
		UserId: model.UserID(order.User),
		OrderID: orderID,
		Status:  model.WaitStatus,
		Message: "the order has been successfully created and is awaiting payment",
	})
	if err != nil {
		return 0, errors.Wrap(err, "Can not notify about order status")
	}

	return model.OrderID(orderID), nil
}
