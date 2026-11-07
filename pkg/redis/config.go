package redis

import "github.com/kelseyhightower/envconfig"

type config struct {
	Address  string `default:"localhost:6379"`
	Password string `default:""`
	DB       int    `default:"0"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
