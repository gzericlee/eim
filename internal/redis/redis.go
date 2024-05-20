package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	RedisEndpoints       []string
	RedisPassword        string
	OfflineMessageExpire time.Duration
	OfflineDeviceExpire  time.Duration
}

type Manager struct {
	redisClient          redis.UniversalClient
	cache                cmap.ConcurrentMap[string, string]
	offlineMessageExpire time.Duration
	offlineDeviceExpire  time.Duration
}

func NewManager(cfg Config) (*Manager, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.RedisEndpoints,
		Password:     cfg.RedisPassword,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MaxRedirects: 8,
		PoolTimeout:  30 * time.Second,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping -> %w", err)
	}

	return &Manager{
		redisClient:          redisClient,
		offlineDeviceExpire:  cfg.OfflineDeviceExpire,
		offlineMessageExpire: cfg.OfflineMessageExpire,
		cache:                cmap.New[string](),
	}, nil
}

func (its *Manager) GetRedisClient() redis.UniversalClient {
	return its.redisClient
}

func (its *Manager) Incr(key string) (int64, error) {
	result, err := its.redisClient.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr -> %w", err)
	}
	return result, nil
}

func (its *Manager) Decr(key string) (int64, error) {
	result, err := its.redisClient.Decr(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis decr -> %w", err)
	}
	return result, nil
}

func (its *Manager) getAllKeys(pattern string) ([]string, error) {
	var errGroup errgroup.Group
	keys := &sync.Map{}

	var scan = func(ctx context.Context, client redis.UniversalClient) error {
		var cursor uint64
		for {
			var keysSlice []string
			var err error
			keysSlice, cursor, err = client.Scan(ctx, cursor, pattern, 1000).Result()
			if err != nil {
				return fmt.Errorf("redis scan -> %w", err)
			}
			for _, key := range keysSlice {
				keys.Store(key, struct{}{})
			}
			if cursor == 0 {
				break
			}
		}
		return nil
	}

	if clusterClient, isOk := its.redisClient.(*redis.ClusterClient); isOk {
		err := clusterClient.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
			errGroup.Go(func() error {
				return scan(ctx, client)
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("redis forEachMaster -> %w", err)
		}
	} else {
		errGroup.Go(func() error {
			return scan(context.Background(), its.redisClient)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf("redis scan error -> %w", err)
	}

	var allKeys []string
	keys.Range(func(key, _ interface{}) bool {
		allKeys = append(allKeys, key.(string))
		return true
	})

	return allKeys, nil
}

func (its *Manager) getAllValues(pattern string) ([]string, error) {
	var errGroup errgroup.Group
	results := &sync.Map{}

	var scan = func(ctx context.Context, client redis.UniversalClient) error {
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = client.Scan(ctx, cursor, pattern, 1000).Result()
			if err != nil {
				return fmt.Errorf("redis scan -> %w", err)
			}
			pipe := client.Pipeline()
			cmds := make([]redis.Cmder, len(keys))
			for i, key := range keys {
				cmds[i] = pipe.Get(ctx, key)
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				return fmt.Errorf("redis pipeline exec -> %w", err)
			}
			for _, cmd := range cmds {
				results.Store(cmd.(*redis.StringCmd).Val(), struct{}{})
			}
			if cursor == 0 {
				break
			}
		}
		return nil
	}

	if clusterClient, isOk := its.redisClient.(*redis.ClusterClient); isOk {
		err := clusterClient.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
			errGroup.Go(func() error {
				return scan(ctx, client)
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("redis forEachMaster -> %w", err)
		}
	} else {
		errGroup.Go(func() error {
			return scan(context.Background(), its.redisClient)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf("redis scan error -> %w", err)
	}

	var allValues []string
	results.Range(func(key, _ interface{}) bool {
		allValues = append(allValues, key.(string))
		return true
	})

	return allValues, nil
}
