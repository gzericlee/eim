package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	batchSize = 100
)

var (
	msgIdsBuffer = make(map[string][]int64)
	bufferLock   = &sync.Mutex{}
)

func (its *Manager) SaveOfflineMessageIds(msgIds []int64, userId, deviceId, bizId string) error {
	bufferLock.Lock()
	defer bufferLock.Unlock()

	key := fmt.Sprintf("%s.offline.messages.%s.%s", userId, deviceId, bizId)
	msgIdsBuffer[key] = append(msgIdsBuffer[key], msgIds...)

	if len(msgIdsBuffer[key]) >= batchSize {
		return its.flushMessageIds(key)
	}

	return nil
}

func (its *Manager) RemoveOfflineMessageIds(msgIds []int64, userId, deviceId, bizId string) error {
	key := fmt.Sprintf("%s.offline.messages.%s.%s", userId, deviceId, bizId)

	ctx := context.Background()

	_, err := its.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, msgId := range msgIds {
			pipe.LRem(ctx, key, 0, msgId)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("redis pipelined lrem(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetOfflineMessagesByBiz(userId, deviceId, bizId string) ([]string, error) {
	key := fmt.Sprintf("%s.offline.messages.%s.%s", userId, deviceId, bizId)

	err := its.flushMessageIds(key)
	if err != nil {
		return nil, fmt.Errorf("flush message ids -> %w", err)
	}

	result, err := its.getAllValues(key)
	if err != nil {
		return nil, fmt.Errorf("rediss lrange(%s) -> %w", key, err)
	}

	return result, nil
}

func (its *Manager) GetOfflineMessagesByDevice(userId, deviceId string) ([]string, error) {
	key := fmt.Sprintf("%s.offline.messages.%s.*", userId, deviceId)

	keys, err := its.getAllKeys(key)
	if err != nil {
		return nil, fmt.Errorf("redis getAllKeys(%s) -> %w", key, err)
	}

	for _, key := range keys {
		err := its.flushMessageIds(key)
		if err != nil {
			return nil, fmt.Errorf("flush message seq ids -> %w", err)
		}
	}

	result, err := its.getAllValues(key)
	if err != nil {
		return nil, fmt.Errorf("rediss lrange(%s) -> %w", key, err)
	}

	return result, nil
}

func (its *Manager) flushMessageIds(key string) error {
	if len(msgIdsBuffer) == 0 {
		return nil
	}

	ctx := context.Background()

	err := its.redisClient.LPush(ctx, key, msgIdsBuffer[key]).Err()
	if err != nil {
		return fmt.Errorf("redis lpush(%s) -> %w", key, err)
	}

	err = its.redisClient.Expire(ctx, key, its.offlineDeviceExpire).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) offline device key expiry -> %w", key, err)
	}

	delete(msgIdsBuffer, key)

	return nil
}
