package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a Redis client configured from environment variables with the given prefix.
func NewClient(envPrefix string) (*redis.Client, error) {

	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}
