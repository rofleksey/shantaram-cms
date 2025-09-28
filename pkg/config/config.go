package config

import (
	"context"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseApiURL   string `yaml:"base_api_url" validate:"required"`
	BaseFrontURL string `yaml:"base_front_url" validate:"required"`
	BaseAdminURL string `yaml:"base_admin_url" validate:"required"`

	Sentry struct {
		DSN string `yaml:"dsn"`
	} `yaml:"sentry"`

	Log struct {
		Telegram struct {
			Token  string `yaml:"token"`
			ChatID string `yaml:"chat_id"`
		} `yaml:"telegram"`
	} `yaml:"log"`

	Telemetry struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"telemetry"`

	DB struct {
		User     string `yaml:"user" validate:"required"`
		Pass     string `yaml:"pass" validate:"required"`
		Host     string `yaml:"host" validate:"required"`
		Database string `yaml:"database" validate:"required"`
	} `yaml:"db"`

	JWT struct {
		Secret string `yaml:"secret" validate:"required"`
	} `yaml:"jwt"`

	Admin struct {
		Password string `yaml:"password" validate:"required"`
	} `yaml:"admin"`

	Telegram struct {
		Token   string   `yaml:"token" validate:"required"`
		ChatIds []string `yaml:"chat_ids" validate:"required"`
	} `yaml:"telegram"`
}

func Load() (*Config, error) {
	span := sentry.StartSpan(context.Background(), "config.load")
	defer span.Finish()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var result Config
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if result.BaseApiURL == "" {
		result.BaseApiURL = "https://api.shantaram-spb.ru"
	}
	if result.BaseFrontURL == "" {
		result.BaseFrontURL = "https://shantaram-spb.ru"
	}
	if result.BaseAdminURL == "" {
		result.BaseAdminURL = "https://admin.shantaram-spb.ru"
	}

	if result.DB.User == "" {
		result.DB.User = "postgres"
	}
	if result.DB.Pass == "" {
		result.DB.Pass = "postgres"
	}
	if result.DB.Host == "" {
		result.DB.Host = "localhost:5432"
	}
	if result.DB.Database == "" {
		result.DB.Database = "shantaram"
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(result); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &result, nil
}
