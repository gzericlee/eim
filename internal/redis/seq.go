package redis

import "context"

func (its *Manager) Incr(key string) (int64, error) {
	return its.redisClient.Incr(context.Background(), key).Result()
}

func (its *Manager) Decr(key string) (int64, error) {
	return its.redisClient.Decr(context.Background(), key).Result()
}
