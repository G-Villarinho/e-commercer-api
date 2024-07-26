package database

import (
	"context"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.RedisAdress,
		Password: config.Env.RedisPassword,
		DB:       config.Env.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
