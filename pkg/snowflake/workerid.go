package snowflake

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"eim/util/log"
)

type workerIdManager struct {
	redisClient     redis.UniversalClient
	maxWorkerId     int64
	minWorkerId     int64
	workerIdKey     string
	usedWorkerIdKey string
	expiration      time.Duration
}

type workerIdConfig struct {
	redisEndpoints []string
	redisPassword  string
	maxWorkerId    int64
	minWorkerId    int64
}

func wewWorkerIdManager(cfg workerIdConfig) (*workerIdManager, error) {
	if cfg.minWorkerId < 0 || cfg.minWorkerId > cfg.maxWorkerId {
		return nil, fmt.Errorf("min worker id must be less than max worker id and greater than 0")
	}
	if cfg.maxWorkerId < 1 || cfg.maxWorkerId >= 1024 {
		return nil, fmt.Errorf("max worker id must be less than 1024 and greater than 1")
	}
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.redisEndpoints,
		Password:     cfg.redisPassword,
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

	return &workerIdManager{
		expiration:      time.Second * 10,
		workerIdKey:     "snowflake.worker.id",
		usedWorkerIdKey: "snowflake.used.worker.id",
		maxWorkerId:     cfg.maxWorkerId,
		minWorkerId:     cfg.minWorkerId,
		redisClient:     redisClient,
	}, nil
}

func (its *workerIdManager) nextId() (int64, error) {
	ctx := context.Background()

	workerId, err := its.redisClient.Incr(ctx, its.workerIdKey).Result()
	if err != nil {
		return -1, fmt.Errorf("redis incr(%s) -> %w", its.workerIdKey, err)
	}

	if workerId < its.minWorkerId || workerId > its.maxWorkerId {
		workerId = its.minWorkerId - 1
		its.redisClient.Set(ctx, its.workerIdKey, workerId, 0)
		return its.nextId()
	}

	key := fmt.Sprintf("%s.%d", its.usedWorkerIdKey, workerId)

	exist, err := its.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("redis exists(%s) -> %w", key, err)
	}
	if exist == 1 {
		return -1, fmt.Errorf("id(%d) is already being used", workerId)
	}

	err = its.redisClient.Set(ctx, key, workerId, its.expiration).Err()
	if err != nil {
		return -1, fmt.Errorf("redis set(%s) -> %w", its.usedWorkerIdKey, err)
	}

	go its.autoRenewExpiration(workerId)

	return workerId, nil
}

func (its *workerIdManager) autoRenewExpiration(workerId int64) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()
	key := fmt.Sprintf("%s.%d", its.usedWorkerIdKey, workerId)

	for {
		select {
		case <-ticker.C:
			ttl, err := its.redisClient.TTL(ctx, key).Result()
			if err != nil {
				log.Error("Error ttl used worker id", zap.Error(err))
				continue
			}

			if ttl <= 5*time.Second {
				err = its.redisClient.Expire(ctx, key, its.expiration).Err()
				if err != nil {
					log.Error("Error to auto expiration used worker id", zap.Error(err))
				}
			}
		}
	}
}
