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

type GroupMessageHandler struct {
	StorageRpc   *storagerpc.Client
	RedisManager *redis.Manager
	Producer     mq.Producer
}

func (its *GroupMessageHandler) HandleMessage(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	msg := &model.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("Error deserializing message", zap.Error(err))
		return nil
	}

	err = toGroup(msg, its.RedisManager, its.Producer)
	if err != nil {
		log.Error("Error dispatching group message to user", zap.String("userId", msg.UserId), zap.Error(err))
		return err
	}

	err = its.StorageRpc.SaveMessage(msg)
	if err != nil {
		log.Error("Error saving message", zap.Error(err))
		return err
	}

	return nil
}
