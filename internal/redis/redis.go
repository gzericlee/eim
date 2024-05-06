package redis

import (
	"context"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	redisClient *redis.ClusterClient
	cache       cmap.ConcurrentMap[string, string]
}

func NewManager(redisEndpoints []string, redisPassword string) (*Manager, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        redisEndpoints,
		Password:     redisPassword,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MaxRedirects: 8,
		PoolTimeout:  30 * time.Second,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &Manager{redisClient: redisClient, cache: cmap.New[string]()}, nil
}

func (its *Manager) GetRedisClient() *redis.ClusterClient {
	return its.redisClient
}

func (its *Manager) getAll(key string) ([]string, error) {
	var locker sync.RWMutex
	var result []string
	var cursor uint64
	err := its.redisClient.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
		iter := client.Scan(ctx, cursor, key, 1000).Iterator()
		for iter.Next(ctx) {
			val, err := client.Get(context.Background(), iter.Val()).Result()
			if err != nil {
				return err
			}
			locker.Lock()
			result = append(result, val)
			locker.Unlock()
		}
		return iter.Err()
	})
	return result, err
}
