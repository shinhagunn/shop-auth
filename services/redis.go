package services

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(address string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: address,
		}),
		ctx: context.Background(),
	}
}

func (c *RedisClient) Set(key string, value interface{}) error {
	_, err := c.client.Set(c.ctx, key, value, 0).Result()
	return err
}

func (c *RedisClient) Get(key string) (interface{}, error) {
	value, err := c.client.Get(c.ctx, key).Result()
	return value, err
}

func (c *RedisClient) Del(key string) error {
	_, err := c.client.Del(c.ctx, key).Result()
	return err
}
