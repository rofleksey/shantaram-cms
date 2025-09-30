package controller

import (
	"context"
	"fmt"
	"net/http"
	"shantaram/app/api"
	"shantaram/app/mapper"

	"github.com/samber/oops"
)

func (s *Server) SetHeaderText(ctx context.Context, request api.SetHeaderTextRequestObject) (api.SetHeaderTextResponseObject, error) {
	if !s.authService.IsAdmin(ctx) {
		return nil, oops.With("status_code", http.StatusUnauthorized).Errorf("Unauthorized")
	}

	if err := s.paramsService.SetHeaderText(ctx, request.Body.Text, request.Body.Deadline); err != nil {
		return nil, fmt.Errorf("SetHeaderText: %w", err)
	}

	return api.SetHeaderText200Response{}, nil
}

func (s *Server) GetParams(ctx context.Context, _ api.GetParamsRequestObject) (api.GetParamsResponseObject, error) {
	params, err := s.queries.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetParams: %w", err)
	}

	return api.GetParams200JSONResponse(mapper.MapParams(params)), nil
}
