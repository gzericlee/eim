package dispatch

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/mq"
	"eim/internal/redis"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
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
		log.Error("Error unmarshal message. Drop it", zap.Error(err))
		return nil
	}

	err = toGroup(msg, its.RedisManager, its.Producer)
	if err != nil {
		return fmt.Errorf("send message to group -> %w", err)
	}

	err = its.StorageRpc.SaveMessage(msg)
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}

	return nil
}
