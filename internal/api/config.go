package api

import (
	"uuid"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the API server settings loaded from environment variables.
type Config struct {
	Port        string `envconfig:"PORT" required:"true"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceId  string `envconfig:"INSTANCE_ID"`
}

// NewConfig loads Config from environment variables.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	if cfg.InstanceId == "" {
		cfg.InstanceId = uuid.New().String()
	}

	return cfg, nil
}
