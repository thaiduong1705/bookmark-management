package redis

import "github.com/redis/go-redis/v9"

type rdsClient struct {
	client *redis.Client
}

var rds rdsClient

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

	rds.client = client

	return client, nil
}

// Close closes the client created by NewClient.
func Close() {
	rds.client.Close()
}
