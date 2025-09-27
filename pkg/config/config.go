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
	BaseApiURL   string `yaml:"baseApiURL" validate:"required"`
	BaseFrontURL string `yaml:"baseFrontURL" validate:"required"`

	Sentry struct {
		DSN              string  `yaml:"dsn"`
		Environment      string  `yaml:"environment"`
		TracesSampleRate float64 `yaml:"traces_sample_rate"`
	} `yaml:"sentry"`

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
}

type Account struct {
	Username     string `yaml:"username" validate:"required"`
	RefreshToken string `yaml:"refresh_token" validate:"required"`
}

func Load() (*Config, error) {
	span := sentry.StartSpan(context.Background(), "config.load")
	defer span.Finish()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		sentry.CaptureException(err)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var result Config
	if err := yaml.Unmarshal(data, &result); err != nil {
		sentry.CaptureException(err)
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if result.BaseApiURL == "" {
		result.BaseApiURL = "https://api.shantaram-spb.ru"
	}
	if result.BaseFrontURL == "" {
		result.BaseFrontURL = "https://shantaram-spb.ru"
	}

	if result.Sentry.TracesSampleRate == 0 {
		result.Sentry.TracesSampleRate = 1.0
	}
	if result.Sentry.Environment == "" {
		result.Sentry.Environment = "production"
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
		sentry.CaptureException(err)
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &result, nil
}
