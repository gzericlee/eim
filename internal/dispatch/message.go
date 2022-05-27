package dispatch

import (
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/nsq/producer"
	"eim/internal/storage"
	"eim/model"
	"eim/proto/pb"
)

type MessageHandler struct {
	StorageRpc *storage.RpcClient
}

func (its *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			msg := &pb.Message{}
			err := proto.Unmarshal(m.Body, msg)
			if err != nil {
				global.Logger.Error("Error deserializing message", zap.Error(err))
				m.Finish()
				return
			}

			//TODO 解析消息分发

			//err = producer.PublishAsync(model.MessageStoreTopic, m.Body)
			err = its.StorageRpc.SaveMessage(msg)
			if err != nil {
				global.Logger.Error("Error saving message", zap.Error(err))
				m.Requeue(-1)
				return
			}

			err = producer.PublishAsync(model.MessageSendTopic, m.Body)
			if err != nil {
				global.Logger.Error("Error publishing message", zap.Error(err))
				m.Requeue(-1)
				return
			}

			m.Finish()
		}
	}(m))

	return nil
}
