package idgenerator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient redis.UniversalClient
var once sync.Once

var workerId int64
var maxWorkerId int64 = 0
var minWorkerId int64 = 0

var workerIdKeyPrefix string = "snowflake.worker.id"

type config struct {
	redisEndpoints []string
	redisPassword  string
	maxWorkerId    int64
	minWorkerId    int64
}

func registry(cfg config) (int64, error) {
	if cfg.maxWorkerId < 0 || cfg.minWorkerId > cfg.maxWorkerId {
		return -2, fmt.Errorf("invalid worker id range")
	}

	once.Do(func() {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.redisEndpoints,
			Password:     cfg.redisPassword,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			MaxRedirects: 8,
			PoolTimeout:  30 * time.Second,
		})
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return -1, fmt.Errorf("ping redis -> %w", err)
	}

	maxWorkerId = cfg.maxWorkerId
	minWorkerId = cfg.minWorkerId

	workerId, err = getNextWorkerId()
	if err != nil {
		return -1, fmt.Errorf("get next worker id -> %w", err)
	}

	return workerId, nil
}

func getNextWorkerId() (int64, error) {
	result, err := redisClient.Incr(ctx, workerIdKeyPrefix).Result()
	if err != nil {
		return -1, fmt.Errorf("redis incr -> %w", err)
	}

	candidateId := int64(result)

	if candidateId < minWorkerId || candidateId > maxWorkerId {
		candidateId = minWorkerId
		redisClient.Set(ctx, workerIdKeyPrefix, candidateId, 0)
	}

	return candidateId, nil
}
