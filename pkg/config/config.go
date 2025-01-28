package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	JWTKey    string `yaml:"jwtKey" validate:"required"`
	AdminPass string `yaml:"adminPass" validate:"required"`
}

func Load() (*Config, error) {
	rawYAML, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var result Config

	err = yaml.Unmarshal(rawYAML, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yaml: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(result); err != nil {
		return nil, fmt.Errorf("failed to validate config.yaml: %w", err)
	}

	return &result, nil
}
