package controller

import (
	"context"
	"shantaram/app/api"
)

func (s *Server) HealthCheck(_ context.Context, _ api.HealthCheckRequestObject) (api.HealthCheckResponseObject, error) {
	return api.HealthCheck200Response{}, nil
}
