package menu

import (
	"context"
	"fmt"
	"net/http"
	"shantaram/app/api"
	"shantaram/app/mapper"
	"shantaram/app/service/pubsub"
	"shantaram/pkg/config"
	"shantaram/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/samber/oops"
)

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

func (s *Service) GetMenu(ctx context.Context) ([]api.Menu, error) {
	menus, err := s.queries.GetMenus(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetMenus: %w", err)
	}

	result := make([]api.Menu, 0, len(menus))
	for _, menu := range menus {
		result = append(result, mapper.MapMenu(menu))
	}

	groups, err := s.queries.GetAllProductGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllProductGroups: %w", err)
	}

	for _, group := range groups {
		for menuIndex := range result {
			if result[menuIndex].Id == group.MenuID {
				result[menuIndex].Groups = append(result[menuIndex].Groups, mapper.MapProductGroup(group))
			}
		}
	}

	products, err := s.queries.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
	}

	for _, product := range products {
		for menuIndex := range result {
			for groupIndex := range result[menuIndex].Groups {
				if result[menuIndex].Groups[groupIndex].Id == product.GroupID {
					result[menuIndex].Groups[groupIndex].Products = append(result[menuIndex].Groups[groupIndex].Products, mapper.MapProduct(product))
				}
			}
		}
	}

	return result, nil
}

func (s *Service) SetMenuOrdering(ctx context.Context, req *api.SetMenuOrderingRequest) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := s.queries.WithTx(tx)

	for index, productGroupID := range req.ProductGroupIds {
		productGroup, err := qtx.GetProductGroupByID(ctx, productGroupID)
		if err != nil {
			return fmt.Errorf("GetProductGroupByID %s: %w", productGroupID, err)
		}

		if productGroup.MenuID != req.MenuId {
			return oops.With("status_code", http.StatusBadRequest).
				Errorf("menuId %s of product group %s does not match menu id %s",
					productGroup.MenuID, productGroup.ID, req.MenuId)
		}

		if err := qtx.UpdateProductGroupIndex(ctx, database.UpdateProductGroupIndexParams{
			ID:    productGroupID,
			Index: int32(index),
		}); err != nil {
			return fmt.Errorf("UpdateProductGroupIndex: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Commit: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) SetProductGroupOrdering(ctx context.Context, req *api.SetProductGroupOrderingRequest) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := s.queries.WithTx(tx)

	for index, productId := range req.ProductIds {
		product, err := qtx.GetProductByID(ctx, productId)
		if err != nil {
			return fmt.Errorf("GetProductByID %s: %w", productId, err)
		}

		if product.GroupID != req.ProductGroupId {
			return oops.With("status_code", http.StatusBadRequest).
				Errorf("productGroupId %s of product %s does not match group id %s",
					product.GroupID, product.ID, req.ProductGroupId)
		}

		if err := qtx.UpdateProductIndex(ctx, database.UpdateProductIndexParams{
			ID:    productId,
			Index: int32(index),
		}); err != nil {
			return fmt.Errorf("UpdateProductIndex: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("DeleteProduct: %w", err)
	}

	s.pubsubService.NotifyOrdersChanged()

	return nil
}

func (s *Service) EditProduct(ctx context.Context, id uuid.UUID, req *api.EditProductRequest) error {
	if err := s.queries.UpdateProduct(ctx, database.UpdateProductParams{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Available:   req.Available,
	}); err != nil {
		return fmt.Errorf("EditProduct: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) DeleteProductGroup(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteProductGroup(ctx, id); err != nil {
		return fmt.Errorf("DeleteProductGroup: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) EditProductGroup(ctx context.Context, id uuid.UUID, req *api.EditProductGroupRequest) error {
	if err := s.queries.UpdateProductGroup(ctx, database.UpdateProductGroupParams{
		ID:    id,
		Title: req.Title,
	}); err != nil {
		return fmt.Errorf("EditProductGroup: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) AddProduct(ctx context.Context, req *api.AddProductRequest) error {
	if err := s.queries.CreateProduct(ctx, database.CreateProductParams{
		ID:          req.Id,
		GroupID:     req.GroupId,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	}); err != nil {
		return fmt.Errorf("AddProduct: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}

func (s *Service) AddProductGroup(ctx context.Context, req *api.AddProductGroupRequest) error {
	if err := s.queries.CreateProductGroup(ctx, database.CreateProductGroupParams{
		ID:     req.Id,
		MenuID: req.MenuId,
		Title:  req.Title,
	}); err != nil {
		return fmt.Errorf("AddProductGroup: %w", err)
	}

	s.pubsubService.NotifyMenuChanged()

	return nil
}
