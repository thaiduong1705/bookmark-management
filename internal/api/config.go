package api

import (
	"github.com/kelseyhightower/envconfig"
)

// Config holds the API server settings loaded from environment variables.
type Config struct {
	Port string `envconfig:"PORT" required:"true"`
}

// NewConfig loads Config from environment variables.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
