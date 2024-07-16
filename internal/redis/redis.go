package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	redisClient redis.UniversalClient

	offlineMessageExpire time.Duration
	offlineDeviceExpire  time.Duration
}

func NewManager(cfg *Config) (*Manager, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.RedisEndpoints,
		Password:     cfg.RedisPassword,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     1000,
		PoolTimeout:  30 * time.Second,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping -> %w", err)
	}

	manager := &Manager{
		redisClient:          redisClient,
		offlineDeviceExpire:  cfg.OfflineDeviceExpire,
		offlineMessageExpire: cfg.OfflineMessageExpire,
	}

	return manager, nil
}

func (its *Manager) GetRedisClient() redis.UniversalClient {
	return its.redisClient
}

func (its *Manager) incr(key string) (int64, error) {
	result, err := its.redisClient.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr -> %w", err)
	}
	return result, nil
}

func (its *Manager) decr(key string) (int64, error) {
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
			keysSlice, cursor, err = client.Scan(ctx, cursor, pattern, 5000).Result()
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
	var mutex sync.Mutex
	var errGroup errgroup.Group
	var allValues []string

	var scan = func(ctx context.Context, client redis.UniversalClient) error {
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = client.Scan(ctx, cursor, pattern, 5000).Result()
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
			mutex.Lock()
			for _, cmd := range cmds {
				allValues = append(allValues, cmd.(*redis.StringCmd).Val())
			}
			mutex.Unlock()
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

	return allValues, nil
}

func (its *Manager) lLenAll(pattern string) (int64, error) {
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
				cmds[i] = pipe.LLen(ctx, key)
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				return fmt.Errorf("redis pipeline exec -> %w", err)
			}
			for _, cmd := range cmds {
				results.Store(cmd.(*redis.IntCmd).Val(), struct{}{})
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
			return 0, fmt.Errorf("redis forEachMaster -> %w", err)
		}
	} else {
		errGroup.Go(func() error {
			return scan(context.Background(), its.redisClient)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return 0, fmt.Errorf("redis scan error -> %w", err)
	}

	var total int64
	results.Range(func(key, _ interface{}) bool {
		total += key.(int64)
		return true
	})

	return total, nil
}

func (its *Manager) lRangeAll(pattern string) (map[string][]string, error) {
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
			cmds := make(map[string]*redis.StringSliceCmd, len(keys))
			for _, key := range keys {
				cmds[key] = pipe.LRange(ctx, key, 0, -1)
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				return fmt.Errorf("redis pipeline exec -> %w", err)
			}
			for key, cmd := range cmds {
				results.Store(key, cmd.Val())
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

	allValues := make(map[string][]string)
	results.Range(func(key, value interface{}) bool {
		allValues[key.(string)] = value.([]string)
		return true
	})

	return allValues, nil
}
