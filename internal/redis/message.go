package redis

import (
	"context"
	"fmt"
	"sync"
)

const (
	batchSize = 100
)

var (
	msgSeqIdsBuffer = make(map[string][]int64)
	bufferLock      = &sync.Mutex{}
)

func (its *Manager) SaveOfflineMessageIds(msgSeqIds []int64, userId, deviceId, bizId string) error {
	bufferLock.Lock()
	defer bufferLock.Unlock()

	key := fmt.Sprintf("%s.offline_message.%s.%s", userId, deviceId, bizId)
	msgSeqIdsBuffer[key] = append(msgSeqIdsBuffer[key], msgSeqIds...)

	if len(msgSeqIdsBuffer[key]) >= batchSize {
		return its.flushMessageIds(key)
	}

	return nil
}

func (its *Manager) flushMessageIds(key string) error {
	if len(msgSeqIdsBuffer) == 0 {
		return nil
	}

	ctx := context.Background()

	err := its.redisClient.LPush(ctx, key, msgSeqIdsBuffer[key]).Err()
	if err != nil {
		return fmt.Errorf("redis lpush(%s) -> %w", key, err)
	}
	err = its.redisClient.Expire(ctx, key, its.offlineDeviceExpire).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) offline device key expiry -> %w", key, err)
	}

	delete(msgSeqIdsBuffer, key)

	return nil
}

func (its *Manager) GetOfflineMessagesByBiz(userId, deviceId, bizId string) ([]string, error) {
	key := fmt.Sprintf("%s.offline_message.%s.%s", userId, deviceId, bizId)
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
	key := fmt.Sprintf("%s.offline_message.%s.*", userId, deviceId)
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
