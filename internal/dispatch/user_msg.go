package dispatch

import (
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/mq"
	"eim/internal/redis"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
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
		log.Error("Error deserializing message", zap.Error(err))
		return nil
	}

	//发送者多端同步
	msg.UserId = msg.FromId
	err = toUser(msg, its.RedisManager, its.Producer)
	if err != nil {
		log.Error("Error dispatching user message to user", zap.String("userId", msg.UserId), zap.Error(err))
		return err
	}

	//接收者多端推送
	msg.UserId = msg.ToId
	err = toUser(msg, its.RedisManager, its.Producer)
	if err != nil {
		log.Error("Error dispatching user message to user", zap.String("userId", msg.UserId), zap.Error(err))
		return err
	}

	err = its.StorageRpc.SaveMessage(msg)
	if err != nil {
		log.Error("Error saving message", zap.Error(err))
		return err
	}

	return nil
}
