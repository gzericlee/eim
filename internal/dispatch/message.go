package dispatch

import (
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"

	"eim/global"
	"eim/internal/nsq/producer"
	"eim/model"
	"eim/proto/pb"
)

type MessageHandler struct {
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
				global.Logger.Warnf("Error deserializeing message: %v", err)
				m.Finish()
				return
			}

			//TODO 解析消息分发

			err = producer.PublishAsync(model.MessageStoreTopic, m.Body)
			if err != nil {
				global.Logger.Warnf("Error publishing message: %v", err)
				m.Requeue(-1)
				return
			}

			err = producer.PublishAsync(model.MessageSendTopic, m.Body)
			if err != nil {
				global.Logger.Warnf("Error publishing message: %v", err)
				m.Requeue(-1)
				return
			}

			m.Finish()
		}
	}(m))

	return nil
}
