package redis

import "github.com/redis/go-redis/v9"

type rdsClient struct {
	client *redis.Client
}

var rds rdsClient

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

func Close() {
	rds.client.Close()
}
