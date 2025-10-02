package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"shantaram/app/api"
	"shantaram/app/mapper"
	"shantaram/app/service/pubsub"
	"shantaram/app/service/telegram"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"shantaram/pkg/telemetry"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/samber/oops"
)

var serviceName = "order"

var maxPositions = 10
var maxPrice = 99999.0
var maxAmount = 10

type Service struct {
	cfg             *config.Config
	dbConn          *pgxpool.Pool
	queries         *database.Queries
	pubsubService   *pubsub.Service
	telegramService *telegram.Service
	tracing         *telemetry.Tracing
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:             do.MustInvoke[*config.Config](di),
		dbConn:          do.MustInvoke[*pgxpool.Pool](di),
		queries:         do.MustInvoke[*database.Queries](di),
		pubsubService:   do.MustInvoke[*pubsub.Service](di),
		telegramService: do.MustInvoke[*telegram.Service](di),
		tracing:         do.MustInvoke[*telemetry.Tracing](di),
	}, nil
}

func (s *Service) CreateOrder(ctx context.Context, req *api.NewOrderRequest) error {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "create")
	defer span.End()

	if len(req.Items) > maxPositions {
		return s.tracing.Error(span, oops.With("status_code", http.StatusBadRequest).New("too many items"))
	}

	var totalPrice float64

	orderItems := make([]api.OrderItem, 0, len(req.Items))
	for _, newItem := range req.Items {
		if newItem.Amount > maxAmount {
			return s.tracing.Error(span, oops.With("status_code", http.StatusBadRequest).New("too many items"))
		}

		item, err := s.mapNewOrderItem(ctx, newItem)
		if err != nil {
			return s.tracing.Error(span, fmt.Errorf("mapNewOrderItem %d: %w", newItem.Id, err))
		}

		orderItems = append(orderItems, item)
		totalPrice += item.Price
	}

	if totalPrice > maxPrice {
		return s.tracing.Error(span, oops.With("status_code", http.StatusBadRequest).New("too many items"))
	}

	dbOrder, err := s.queries.CreateOrder(ctx, database.CreateOrderParams{
		ID:            req.Id,
		TableID:       nil,
		ClientName:    req.Name,
		ClientComment: req.Comment,
		Status:        api.OrderStatusOpen,
		Seen:          false,
		Items:         orderItems,
	})
	if err != nil {
		return s.tracing.Error(span, fmt.Errorf("CreateOrder: %w", err))
	}

	msg := mapper.OrderToNotificationText(dbOrder)
	go s.telegramService.Notify(msg)

	s.pubsubService.NotifyOrdersChanged()
	s.tracing.Success(span)

	return nil
}

func (s *Service) SetStatus(ctx context.Context, id uuid.UUID, status api.OrderStatus) error {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "set_status")
	defer span.End()

	if err := s.queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     id,
		Status: status,
	}); err != nil {
		return s.tracing.Error(span, fmt.Errorf("UpdateOrderStatus: %w", err))
	}

	s.pubsubService.NotifyOrdersChanged()
	s.tracing.Success(span)

	return nil
}

func (s *Service) DeleteOrderByID(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "set_status")
	defer span.End()

	if err := s.queries.DeleteOrder(ctx, id); err != nil {
		return s.tracing.Error(span, fmt.Errorf("DeleteOrder: %w", err))
	}

	s.pubsubService.NotifyOrdersChanged()
	s.tracing.Success(span)

	return nil
}

func (s *Service) GetOrderByID(ctx context.Context, id uuid.UUID) (database.Order, error) {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "get_order_by_id")
	defer span.End()

	order, err := s.queries.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Order{}, s.tracing.Error(span, oops.With("status_code", http.StatusNotFound).Errorf("order not found"))
		}

		return database.Order{}, s.tracing.Error(span, fmt.Errorf("GetOrderByID: %w", err))
	}

	s.tracing.Success(span)

	return order, nil
}

func (s *Service) GetOrdersPaginated(ctx context.Context, offset, limit int) ([]database.Order, int64, error) {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "get_orders_paginated")
	defer span.End()

	orders, err := s.queries.GetOrdersPaginated(ctx, database.GetOrdersPaginatedParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, 0, s.tracing.Error(span, fmt.Errorf("GetOrdersPaginated: %w", err))
	}

	totalCount, err := s.queries.CountOrders(ctx)
	if err != nil {
		return nil, 0, s.tracing.Error(span, fmt.Errorf("CountOrders: %w", err))
	}

	s.tracing.Success(span)

	return orders, totalCount, nil
}

func (s *Service) MarkOrderSeen(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "mark_order_seen")
	defer span.End()

	if err := s.queries.SetOrderSeen(ctx, id); err != nil {
		return s.tracing.Error(span, fmt.Errorf("SetOrderSeen: %w", err))
	}

	s.tracing.Success(span)

	return nil
}
