package redis

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"
)

const (
	batchSize                = 100
	offlineMessagesKeyFormat = "offline:messages:%s:%s"
)

func (its *Manager) checkProcessExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		for item := range its.msgsBuffer.IterBuffered() {
			err := its.flushMessages(item.Key)
			if err != nil {
				log.Error("Error flushing redis offline messages", zap.Error(err))
				continue
			}
			log.Info("Flush redis offline messages", zap.String("key", item.Key), zap.Any("count", len(item.Val)))
		}
		os.Exit(0)
	}()
}

func (its *Manager) SaveOfflineMessages(msgs []*model.Message, userId, deviceId string) error {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	its.msgsBuffer.Upsert(key, msgs, func(exist bool, valueInMap, newValue []*model.Message) []*model.Message {
		if !exist {
			return newValue
		}
		return append(valueInMap, newValue...)
	})

	if allMsgs, exist := its.msgsBuffer.Get(key); exist {
		if len(allMsgs) >= batchSize {
			return its.flushMessages(key)
		}
	}

	return nil
}

func (its *Manager) RemoveOfflineMessages(msgIds []string, userId, deviceId string) error {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	err := its.flushMessages(key)
	if err != nil {
		return fmt.Errorf("flush messages -> %w", err)
	}

	err = its.redisClient.HDel(context.Background(), key, msgIds...).Err()
	if err != nil {
		return fmt.Errorf("redis pipelined lrem(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	err := its.flushMessages(key)
	if err != nil {
		return 0, fmt.Errorf("flush messages -> %w", err)
	}

	total, err := its.redisClient.HLen(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis hlen(%s) -> %w", key, err)
	}

	return total, nil
}

func (its *Manager) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	err := its.flushMessages(key)
	if err != nil {
		return nil, fmt.Errorf("flush messages -> %w", err)
	}

	result, err := its.redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange(%s) -> %w", key, err)
	}

	var msgs []*model.Message
	for _, msgStr := range result {
		msg := &model.Message{}
		err := proto.Unmarshal([]byte(msgStr), msg)
		if err != nil {
			return nil, fmt.Errorf("unmarshal message -> %w", err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (its *Manager) flushMessages(key string) error {
	msgs, exist := its.msgsBuffer.Get(key)
	if !exist || msgs == nil || len(msgs) == 0 {
		return nil
	}

	ctx := context.Background()

	data := make(map[string]string, len(msgs))
	for _, msg := range msgs {
		bytes, err := proto.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal message -> %w", err)
		}
		data[strconv.FormatInt(msg.MsgId, 10)] = string(bytes)
	}

	err := its.redisClient.HMSet(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("redis hmset(%s) -> %w", key, err)
	}

	err = its.redisClient.Expire(ctx, key, its.offlineDeviceExpire).Err()
	if err != nil {
		return fmt.Errorf("redis expire(%s) messages -> %w", key, err)
	}

	its.msgsBuffer.Remove(key)

	return nil
}
