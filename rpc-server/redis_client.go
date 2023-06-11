package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func (c *RedisClient) InitClient(ctx context.Context, address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	if err := r.Ping(ctx).Err(); err != nil {
		return err
	}

	c.client = r
	return nil
}
