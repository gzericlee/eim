package redis

import (
	"context"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	redisClient redis.UniversalClient
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

func (its *Manager) GetRedisClient() redis.UniversalClient {
	return its.redisClient
}

func (its *Manager) getAll(key string, limit int64) ([]string, error) {
	var cursor uint64
	var wg sync.WaitGroup
	var mu sync.Mutex
	var scanErr error

	var scan = func(ctx context.Context, client redis.UniversalClient) ([]string, error) {
		var keys []string
		var result []string
		iter := client.Scan(ctx, cursor, key, limit).Iterator()
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}
		pipe := client.Pipeline()
		cmds := make([]redis.Cmder, len(keys))
		for i, key := range keys {
			cmds[i] = pipe.Get(ctx, key)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
		for _, cmd := range cmds {
			result = append(result, cmd.(*redis.StringCmd).Val())
		}
		return result, nil
	}

	results := make([]string, 0)
	if clusterClient, isOk := its.redisClient.(*redis.ClusterClient); isOk {
		clusterClient.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := scan(ctx, client)
				if err != nil {
					scanErr = err
					return
				}
				mu.Lock()
				results = append(results, res...)
				mu.Unlock()
			}()
			return nil
		})
		wg.Wait()
		if scanErr != nil {
			return nil, scanErr
		}
	} else {
		res, err := scan(context.Background(), its.redisClient)
		if err != nil {
			return nil, err
		}
		results = append(results, res...)
	}

	return results, nil
}
