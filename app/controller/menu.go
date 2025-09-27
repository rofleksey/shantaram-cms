package controller

import (
	"context"
	"net/http"
	"shantaram/app/api"

	"github.com/samber/oops"
)

func (s *Server) GetMenu(ctx context.Context, _ api.GetMenuRequestObject) (api.GetMenuResponseObject, error) {
	menus, err := s.menuService.GetMenu(ctx)
	if err != nil {
		return nil, err
	}

	return api.GetMenu200JSONResponse{
		Menus: menus,
	}, nil
}

func (s *Server) SetMenuOrdering(ctx context.Context, req api.SetMenuOrderingRequestObject) (api.SetMenuOrderingResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.SetMenuOrdering(ctx, req.Body); err != nil {
		return nil, err
	}

	return api.SetMenuOrdering200Response{}, nil
}

func (s *Server) SetProductGroupOrdering(ctx context.Context, req api.SetProductGroupOrderingRequestObject) (api.SetProductGroupOrderingResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.SetProductGroupOrdering(ctx, req.Body); err != nil {
		return nil, err
	}

	return api.SetProductGroupOrdering200Response{}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req api.DeleteProductRequestObject) (api.DeleteProductResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.DeleteProduct(ctx, req.ProductId); err != nil {
		return nil, err
	}

	return api.DeleteProduct200Response{}, nil
}

func (s *Server) EditProduct(ctx context.Context, req api.EditProductRequestObject) (api.EditProductResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.EditProduct(ctx, req.ProductId, req.Body); err != nil {
		return nil, err
	}

	return api.EditProduct200Response{}, nil
}

func (s *Server) DeleteProductGroup(ctx context.Context, req api.DeleteProductGroupRequestObject) (api.DeleteProductGroupResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.DeleteProductGroup(ctx, req.ProductGroupId); err != nil {
		return nil, err
	}

	return api.DeleteProductGroup200Response{}, nil
}

func (s *Server) EditProductGroup(ctx context.Context, req api.EditProductGroupRequestObject) (api.EditProductGroupResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.EditProductGroup(ctx, req.ProductGroupId, req.Body); err != nil {
		return nil, err
	}

	return api.EditProductGroup200Response{}, nil
}

func (s *Server) AddProduct(ctx context.Context, req api.AddProductRequestObject) (api.AddProductResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.AddProduct(ctx, req.Body); err != nil {
		return nil, err
	}

	return api.AddProduct200Response{}, nil
}

func (s *Server) AddProductGroup(ctx context.Context, req api.AddProductGroupRequestObject) (api.AddProductGroupResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.menuService.AddProductGroup(ctx, req.Body); err != nil {
		return nil, err
	}

	return api.AddProductGroup200Response{}, nil
}
