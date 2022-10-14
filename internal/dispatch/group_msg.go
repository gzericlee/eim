package dispatch

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/internal/pool"
	"eim/internal/redis"
	storage_rpc "eim/internal/storage/rpc"
	"eim/internal/types"
	"eim/pkg/log"
)

type GroupMessageHandler struct {
	StorageRpc *storage_rpc.Client
}

func (its *GroupMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	msg := &types.Message{}
	err := msg.Deserialize(m.Body)
	if err != nil {
		log.Error("Error deserializing message", zap.Error(err))
		m.Finish()
		return nil
	}

	pool.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			err = its.StorageRpc.SaveMessage(msg)
			if err != nil {
				log.Error("Error saving message", zap.Error(err))
				m.Requeue(-1)
				return
			}
		}
	}(m))

	//发送者多端同步
	msg.UserId = msg.FromId
	err = toUser(msg)
	if err != nil {
		log.Error("Error dispatching message to user", zap.String("userId", msg.UserId), zap.Error(err))
		return err
	}

	//接收者多端推送
	members, err := redis.GetGroupMembers(msg.ToId)
	if err != nil {
		log.Error("Error getting Group members", zap.String("groupId", msg.ToId), zap.Error(err))
		return err
	}
	for _, userId := range members {
		msg.UserId = userId
		err = toUser(msg)
		if err != nil {
			log.Error("Error dispatching message to user", zap.String("userId", msg.UserId), zap.Error(err))
			return err
		}
	}

	m.Finish()

	return nil
}
