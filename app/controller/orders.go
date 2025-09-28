package controller

import (
	"context"
	"fmt"
	"net/http"
	"shantaram/app/api"
	"shantaram/app/mapper"

	"github.com/elliotchance/pie/v2"
	"github.com/samber/oops"
)

func (s *Server) CreateOrder(ctx context.Context, req api.CreateOrderRequestObject) (api.CreateOrderResponseObject, error) {
	if !s.limitsService.AllowIpRpm(ctx, "create_order", 5) {
		return nil, oops.With("status_code", http.StatusTooManyRequests).New("Too many requests")
	}

	if err := s.orderService.CreateOrder(ctx, req.Body); err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	return api.CreateOrder200Response{}, nil
}

func (s *Server) SetOrderStatus(ctx context.Context, request api.SetOrderStatusRequestObject) (api.SetOrderStatusResponseObject, error) {
	if err := s.orderService.SetStatus(ctx, request.Body.Id, request.Body.Status); err != nil {
		return nil, fmt.Errorf("SetStatus: %w", err)
	}

	return api.SetOrderStatus200Response{}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req api.DeleteOrderRequestObject) (api.DeleteOrderResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.orderService.DeleteOrderByID(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("DeleteOrderByID: %w", err)
	}

	return api.DeleteOrder200Response{}, nil
}

func (s *Server) GetOrder(ctx context.Context, req api.GetOrderRequestObject) (api.GetOrderResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	order, err := s.orderService.GetOrderByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("GetOrderByID: %w", err)
	}

	return api.GetOrder200JSONResponse(mapper.MapOrder(order)), nil
}

func (s *Server) GetOrders(ctx context.Context, req api.GetOrdersRequestObject) (api.GetOrdersResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	offset := 0
	limit := 10

	if req.Params.Offset != nil {
		offset = max(*req.Params.Offset, 0)
	}
	if req.Params.Limit != nil {
		limit = min(max(*req.Params.Limit, 0), 100)
	}

	orders, totalCount, err := s.orderService.GetOrdersPaginated(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("GetOrdersPaginated: %w", err)
	}

	return api.GetOrders200JSONResponse{
		Data:       pie.Map(orders, mapper.MapOrder),
		TotalCount: int(totalCount),
	}, nil
}

func (s *Server) MarkOrderSeen(ctx context.Context, req api.MarkOrderSeenRequestObject) (api.MarkOrderSeenResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.orderService.MarkOrderSeen(ctx, req.Body.Id); err != nil {
		return nil, fmt.Errorf("MarkOrderSeen: %w", err)
	}

	return api.MarkOrderSeen200Response{}, nil
}
