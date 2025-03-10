package counter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCounter struct {
	client *redis.Client
}

func NewRedisCounter(client *redis.Client) *RedisCounter {
	return &RedisCounter{client: client}
}

func (r *RedisCounter) NextID() (int64, error) {
	return r.client.Incr(context.Background(), "url_shortener_counter").Result()
}
