package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"shantaram-cms/pkg/config"
	"shantaram-cms/pkg/util"
	"time"
)

type Auth struct {
	cfg *config.Config
}

func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		cfg: cfg,
	}
}

func (s *Auth) AuthAdmin(pass string) (string, error) {
	if pass != s.cfg.AdminPass {
		return "", util.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{}

	claims["id"] = "admin"
	claims["expires"] = time.Now().Add(time.Hour * 24 * 365 * 10).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	result, err := token.SignedString([]byte(s.cfg.JWTKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return result, nil
}
