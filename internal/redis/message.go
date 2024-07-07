package redis

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
)

const (
	offlineMessagesKeyFormat = "messages:offline:%s:%s"
)

func (its *Manager) SaveOfflineMessages(msgs []*model.Message, userId, deviceId string) error {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

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

	return nil
}

func (its *Manager) RemoveOfflineMessages(msgIds []string, userId, deviceId string) error {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	err := its.redisClient.HDel(context.Background(), key, msgIds...).Err()
	if err != nil {
		return fmt.Errorf("redis pipelined lrem(%s) -> %w", key, err)
	}

	return nil
}

func (its *Manager) GetOfflineMessagesCount(userId, deviceId string) (int64, error) {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

	total, err := its.redisClient.HLen(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis hlen(%s) -> %w", key, err)
	}

	return total, nil
}

func (its *Manager) GetOfflineMessages(userId, deviceId string) ([]*model.Message, error) {
	key := fmt.Sprintf(offlineMessagesKeyFormat, userId, deviceId)

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

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].SeqId < msgs[j].SeqId
	})

	return msgs, nil
}
