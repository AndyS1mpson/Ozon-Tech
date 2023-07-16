// Changing order status to paid
package domain

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

// Mark the order as paid
func (s *Service) OrderPayed(ctx context.Context, orderID model.OrderID) error {
	order, err := s.order.GetOrder(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "can not get order")
	}

	_, err = s.stock.WriteOffOrderItems(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "try to write off order items from warehouses")
	}

	err = s.order.PayOrder(ctx, orderID)
	if err != nil {
		return errors.Wrap(err, "try to paid order")
	}

	err = s.notifier.SendMessage(OrderStatusNotification{
		UserId: model.UserID(order.User),
		OrderID: orderID,
		Status:  model.PaidStatus,
		Message: "the order has been successfully paid, we are collecting the goods",
	})
	if err != nil {
		return errors.Wrap(err, "Can not notify about order status")
	}

	return nil
}
