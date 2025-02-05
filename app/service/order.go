package service

import (
	"fmt"
	"shantaram-cms/pkg/database"
)

type Order struct {
	db *database.Database
}

func NewOrder(db *database.Database) *Order {
	return &Order{
		db: db,
	}
}

func (s *Order) GetByID(id uint64) (*database.Order, error) {
	order, err := s.db.GetOrderByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id %d: %v", id, err)
	}

	return order, nil
}

func (s *Order) GetPaginated(offset, limit int) (*database.DataPage[database.Order], error) {
	data, err := s.db.GetOrdersPaginated(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}

	return data, nil
}

func (s *Order) Delete(id uint64) error {
	if err := s.db.DeleteOrder(id); err != nil {
		return fmt.Errorf("failed to delete order by id %d: %v", id, err)
	}

	return nil
}

func (s *Order) Insert(order *database.Order) error {
	if err := s.db.InsertOrder(order); err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	return nil
}

func (s *Order) Update(order *database.Order) error {
	if err := s.db.UpdateOrder(order); err != nil {
		return fmt.Errorf("failed to update order by id %d: %v", order.ID, err)
	}

	return nil
}
