package controller

import (
	"context"
	"net/http"
	"shantaram/app/api"

	"github.com/samber/oops"
)

func (s *Server) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	if !s.limitsService.AllowIpRpm(ctx, "login", 5) {
		return nil, oops.With("status_code", http.StatusTooManyRequests).New("Too many requests")
	}

	token, err := s.authService.Login(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		return nil, err
	}

	return api.Login200JSONResponse{Token: token}, nil
}
