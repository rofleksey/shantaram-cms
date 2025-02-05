package service

import (
	"encoding/json"
	"fmt"
	"shantaram-cms/app/dao"
	"shantaram-cms/pkg/database"
	"time"
)

type Order struct {
	db             *database.Database
	captchaService *Captcha
}

func NewOrder(
	db *database.Database,
	captchaService *Captcha,
) *Order {
	return &Order{
		db:             db,
		captchaService: captchaService,
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

func (s *Order) Insert(req dao.NewOrderRequest) error {
	captchaOk := s.captchaService.Verify(req.CaptchaID, req.CaptchaAnswer)
	if !captchaOk {
		return fmt.Errorf("invalid captcha answer")
	}

	bytez, err := s.db.GetGeneralByIDData("menu")
	if err != nil {
		return fmt.Errorf("failed to get menu: %v", err)
	}

	var settings database.MenuSettings

	if err := json.Unmarshal(bytez, &settings); err != nil {
		return fmt.Errorf("failed to unmarshal menu: %v", err)
	}

	orderItems := make([]database.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		menuItem := settings.FindProduct(item.ID)

		if menuItem == nil {
			return fmt.Errorf("продукт с ID = %s не найден в меню", item.ID)
		}

		orderItems = append(orderItems, database.OrderItem{
			ID:     item.ID,
			Title:  menuItem.Title,
			Price:  menuItem.Price,
			Amount: item.Amount,
		})
	}

	order := database.Order{
		Created: time.Now(),
		Status:  database.OrderStatusOpen,
		Name:    req.Name,
		Phone:   req.Phone,
		Comment: req.Comment,
		Items:   orderItems,
	}

	if err := s.db.InsertOrder(&order); err != nil {
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
