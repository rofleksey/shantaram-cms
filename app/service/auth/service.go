package auth

import (
	"context"
	_ "embed"
	"fmt"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do"
)

type Service struct {
	cfg       *config.Config
	queries   *database.Queries
	startTime time.Time
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:       do.MustInvoke[*config.Config](di),
		queries:   do.MustInvoke[*database.Queries](di),
		startTime: time.Now(),
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

func (s *Service) Login(user, pass string) (string, error) {
	user = strings.TrimSpace(user)

	if user != "admin" || pass != s.cfg.Admin.Password {
		return "", fmt.Errorf("invalid username or password")
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
		"sub": "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}
