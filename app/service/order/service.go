package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"shantaram/app/api"
	"shantaram/app/service/pubsub"
	"shantaram/pkg/config"
	"shantaram/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/samber/oops"
)

var maxPositions = 10
var maxPrice float32 = 99999
var maxAmount = 10

type Service struct {
	cfg           *config.Config
	dbConn        *pgxpool.Pool
	queries       *database.Queries
	pubsubService *pubsub.Service
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:           do.MustInvoke[*config.Config](di),
		dbConn:        do.MustInvoke[*pgxpool.Pool](di),
		queries:       do.MustInvoke[*database.Queries](di),
		pubsubService: do.MustInvoke[*pubsub.Service](di),
	}, nil
}

func (s *Service) CreateOrder(ctx context.Context, req *api.NewOrderRequest) error {
	if len(req.Items) > maxPositions {
		return oops.With("status_code", http.StatusBadRequest).New("too many items")
	}

	var totalPrice float32

	orderItems := make([]api.OrderItem, 0, len(req.Items))
	for _, newItem := range req.Items {
		if newItem.Amount > maxAmount {
			return oops.With("status_code", http.StatusBadRequest).New("too many items")
		}

		item, err := s.mapNewOrderItem(ctx, newItem)
		if err != nil {
			return fmt.Errorf("mapNewOrderItem %d: %w", newItem.Id, err)
		}

		orderItems = append(orderItems, item)
		totalPrice += item.Price
	}

	if totalPrice > maxPrice {
		return oops.With("status_code", http.StatusBadRequest).New("too many items")
	}

	if err := s.queries.CreateOrder(ctx, database.CreateOrderParams{
		ID:            req.Id,
		TableID:       nil,
		ClientName:    req.Name,
		ClientPhone:   req.Phone,
		ClientComment: req.Comment,
		Status:        api.OrderStatusOpen,
		Seen:          false,
		Items:         orderItems,
	}); err != nil {
		return fmt.Errorf("CreateOrder: %w", err)
	}

	// TODO: send message to telegram
	s.pubsubService.NotifyOrdersChanged()

	return nil
}

func (s *Service) SetStatus(ctx context.Context, id uuid.UUID, status api.OrderStatus) error {
	if err := s.queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     id,
		Status: status,
	}); err != nil {
		return fmt.Errorf("UpdateOrderStatus: %w", err)
	}

	return nil
}

func (s *Service) DeleteOrderByID(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteOrder(ctx, id); err != nil {
		return fmt.Errorf("DeleteOrder: %w", err)
	}

	s.pubsubService.NotifyOrdersChanged()

	return nil
}

func (s *Service) GetOrderByID(ctx context.Context, id uuid.UUID) (database.Order, error) {
	order, err := s.queries.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Order{}, oops.With("status_code", http.StatusNotFound).Errorf("order not found")
		}

		return database.Order{}, err
	}

	return order, nil
}

func (s *Service) GetOrdersPaginated(ctx context.Context, offset, limit int) ([]database.Order, int64, error) {
	orders, err := s.queries.GetOrdersPaginated(ctx, database.GetOrdersPaginatedParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("GetOrdersPaginated: %w", err)
	}

	totalCount, err := s.queries.CountOrders(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("CountOrders: %w", err)
	}

	return orders, totalCount, nil
}
