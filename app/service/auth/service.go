package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"shantaram/pkg/telemetry"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do"
	"go.opentelemetry.io/otel/attribute"
)

var serviceName = "auth"

var ErrInvalidCredentials = errors.New("invalid username or password")

type Service struct {
	cfg     *config.Config
	queries *database.Queries
	tracing *telemetry.Tracing
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:     do.MustInvoke[*config.Config](di),
		queries: do.MustInvoke[*database.Queries](di),
		tracing: do.MustInvoke[*telemetry.Tracing](di),
	}, nil
}

func (s *Service) IsAdmin(ctx context.Context) bool {
	isAdmin, _ := ctx.Value("admin").(bool)

	return isAdmin
}

func (s *Service) IsAdminLocals(getter func(key string, value ...interface{}) interface{}) bool {
	isAdmin, _ := getter("admin", false).(bool)

	return isAdmin
}

func (s *Service) Login(ctx context.Context, user, pass string) (string, error) {
	ctx, span := s.tracing.StartServiceSpan(ctx, serviceName, "login")
	defer span.End()

	user = strings.TrimSpace(user)
	span.SetAttributes(attribute.String("username", user))

	success := user == "admin" && pass == s.cfg.Admin.Password
	if !success {
		return "", s.tracing.Error(span, ErrInvalidCredentials)
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
		"sub": "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", s.tracing.Error(span, fmt.Errorf("failed to sign token: %w", err))
	}

	s.tracing.Success(span)

	return tokenStr, nil
}
