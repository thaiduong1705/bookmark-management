package redis

import "github.com/kelseyhightower/envconfig"

// config holds the Redis connection settings.
type config struct {
	Address  string `default:"localhost:6379"`
	Password string `default:""`
	DB       int    `default:"0"`
}

// newConfig loads config from environment variables with the given prefix.
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
