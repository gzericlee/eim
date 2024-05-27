package redis

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"eim/util/log"
)

const (
	batchSize = 100
)

func (its *Manager) checkProcessExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		for item := range its.msgIdsBuffer.IterBuffered() {
			err := its.flushMessageIds(item.Key)
			if err != nil {
				log.Error("Error flushing message ids", zap.Error(err))
				continue
			}
			log.Info("Flush message ids", zap.String("key", item.Key), zap.Any("count", len(item.Val)))
		}
		os.Exit(0)
	}()
}

func (its *Manager) SaveOfflineMessageIds(msgIds []interface{}, userId, deviceId string) error {
	key := fmt.Sprintf("%s:offline:messages:%s", userId, deviceId)

	its.msgIdsBuffer.Upsert(key, msgIds, func(exist bool, valueInMap, newValue []interface{}) []interface{} {
		if !exist {
			return newValue
		}
		return append(valueInMap, newValue...)
	})

	if msgIds, exist := its.msgIdsBuffer.Get(key); exist {
		if len(msgIds) >= batchSize {
			return its.flushMessageIds(key)
		}
	}

	return nil
}

func (its *Manager) RemoveOfflineMessageIds(msgIds []interface{}, userId, deviceId string) error {
	key := fmt.Sprintf("%s:offline:messages:%s", userId, deviceId)

	err := its.flushMessageIds(key)
	if err != nil {
		return fmt.Errorf("flush message ids -> %w", err)
	}

	ctx := context.Background()

	_, err = its.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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

func (its *Manager) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	key := fmt.Sprintf("%s:offline:messages:%s", userId, deviceId)

	total, err := its.lLenAll(key)
	if err != nil {
		return 0, fmt.Errorf("redis foreachmaster llen(%s) -> %w", key, err)
	}

	return total, nil
}

func (its *Manager) GetOfflineMessages(userId, deviceId string) (map[string][]string, error) {
	key := fmt.Sprintf("%s:offline:messages:%s", userId, deviceId)

	err := its.flushMessageIds(key)
	if err != nil {
		return nil, fmt.Errorf("flush message seq ids -> %w", err)
	}

	result, err := its.lRangeAll(key)
	if err != nil {
		return nil, fmt.Errorf("rediss lrange(%s) -> %w", key, err)
	}

	return result, nil
}

func (its *Manager) flushMessageIds(key string) error {
	msgIds, exist := its.msgIdsBuffer.Get(key)
	if !exist || msgIds == nil || len(msgIds) == 0 {
		return nil
	}

	ctx := context.Background()

	interfaceMsgIds := make([]interface{}, len(msgIds))
	for i, v := range msgIds {
		interfaceMsgIds[i] = v
	}

	err := its.redisClient.LPush(ctx, key, interfaceMsgIds...).Err()
	if err != nil {
		return fmt.Errorf("redis lpush(%s) -> %w", key, err)
	}

	err = its.redisClient.Expire(ctx, key, its.offlineDeviceExpire).Err()
	if err != nil {
		return fmt.Errorf("redis set(%s) offline device key expiry -> %w", key, err)
	}

	its.msgIdsBuffer.Remove(key)

	return nil
}
