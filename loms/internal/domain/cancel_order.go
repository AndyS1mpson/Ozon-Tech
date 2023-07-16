// Cancel order
package domain

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

// Cancel the order and removes the reserve from all items
func (s *Service) CancelOrder(ctx context.Context, orderID model.OrderID) error {
	order, err := s.order.GetOrder(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "can not get order")
	}

	for _, item := range order.Items {
		err := s.stock.Unreserve(ctx, orderID, model.SKU(item.SKU))
		if err != nil {
			return errors.Wrap(err, "can not unreserve item")
		}
	}

	err = s.order.CancelOrder(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "try to cancel order")
	}

	err = s.notifier.SendMessage(OrderStatusNotification{
		UserId: model.UserID(order.User),
		OrderID: orderID,
		Status:  model.CanceledStatus,
		Message: "order canceled",
	})
	if err != nil {
		return errors.Wrap(err, "Can not notify about order status")
	}

	return nil
}
