package dispatch

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"eim/internal/model"
	"eim/internal/mq"
	"eim/internal/redis"
	storagerpc "eim/internal/storage/rpc"
)

type UserMessageHandler struct {
	StorageRpc   *storagerpc.Client
	RedisManager *redis.Manager
	Producer     mq.Producer
}

func (its *UserMessageHandler) HandleMessage(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	msg := &model.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	err = its.StorageRpc.SaveMessage(msg)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	//发送者多端同步
	msg.UserId = msg.FromId
	err = toUser(msg, its.RedisManager, its.Producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to send user -> %w", err)
	}

	//接收者多端推送
	msg.UserId = msg.ToId
	err = toUser(msg, its.RedisManager, its.Producer)
	if err != nil {
		return fmt.Errorf("dispatch user message to receive user -> %w", err)
	}

	return nil
}
