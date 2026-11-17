package api

import (
	"github.com/google/uuid"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the API server settings loaded from environment variables.
type Config struct {
	Port        string `envconfig:"APP_PORT" required:"true"`
	ServiceName string `envconfig:"APP_SERVICE_NAME" default:"bookmark_service"`
	InstanceId  string `envconfig:"APP_INSTANCE_ID"`
	BasePath    string `envconfig:"APP_BASE_PATH" default:"/"`
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
