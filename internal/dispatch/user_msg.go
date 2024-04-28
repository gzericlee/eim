package dispatch

import (
	"github.com/nsqio/go-nsq"
	"github.com/panjf2000/ants"
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
	TaskPool     *ants.Pool
}

func (its *UserMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := its.TaskPool.Submit(func(m *nsq.Message) func() {
		return func() {
			msg := &model.Message{}
			err := msg.Deserialize(m.Body)
			if err != nil {
				log.Error("Error deserializing message", zap.Error(err))
				m.Finish()
				return
			}

			err = its.StorageRpc.SaveMessage(msg)
			if err != nil {
				log.Error("Error saving message", zap.Error(err))
				m.Requeue(-1)
				return
			}

			//发送者多端同步
			msg.UserId = msg.FromId
			err = toUser(msg, its.RedisManager, its.Producer)
			if err != nil {
				log.Error("Error dispatching user message to user", zap.String("userId", msg.UserId), zap.Error(err))
				return
			}

			//接收者多端推送
			msg.UserId = msg.ToId
			err = toUser(msg, its.RedisManager, its.Producer)
			if err != nil {
				log.Error("Error dispatching user message to user", zap.String("userId", msg.UserId), zap.Error(err))
				return
			}

			m.Finish()
		}
	}(m))
	if err != nil {
		return err
	}

	return nil
}
