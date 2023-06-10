// Creating a custom order
package domain

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
)

// Create user order
func (s *Service) CreateOrder(ctx context.Context, order model.Order) (model.OrderID, error) {

	// Create order
	orderID, err := s.order.CreateOrder(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %s", err)
	}

	// Reserve items from order
	for _, v := range order.Items {
		stocks, err := s.stock.GetAvailableStocks(ctx, model.SKU(v.SKU))
		if err != nil {
			s.order.FailOrder(ctx, orderID)
			return 0, fmt.Errorf("no available item stocks: %s", err)
		}

		fmt.Printf("%+v\n", stocks)

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
				return 0, fmt.Errorf("failed to update stock count: %s", err)
			}

			if remainingQuantity == 0 {
				break
			}
		}

		if remainingQuantity != 0 {
			s.order.FailOrder(ctx, orderID)
			return 0, fmt.Errorf("not enough available stocks")
		}
	}

	err = s.order.AwaitPaymentOrder(ctx, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to await payment order: %s", err)
	}

	return model.OrderID(orderID), nil
}
