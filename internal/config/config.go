package config

import (
	"github.com/google/uuid"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        string `envconfig:"PORT" required:"true"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.NewString()
	}

	return cfg, nil
}
