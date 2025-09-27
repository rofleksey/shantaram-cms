package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shantaram/app/api"
)

func (s *Service) mapNewOrderItem(ctx context.Context, item api.NewOrderItem) (api.OrderItem, error) {
	product, err := s.queries.GetProductByID(ctx, item.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.OrderItem{}, fmt.Errorf("product %s not found: %w", item.Id, err)
		}

		return api.OrderItem{}, fmt.Errorf("GetProductByID: %w", err)
	}

	return api.OrderItem{
		Amount: item.Amount,
		Id:     item.Id,
		Price:  product.Price,
		Title:  product.Title,
	}, nil
}
