package storage

import (
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/model"
	"eim/proto/pb"
)

type MessageHandler struct{}

func (its *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			msg := &pb.Message{}
			err := proto.Unmarshal(m.Body, msg)
			if err != nil {
				m.Finish()
				return
			}

			dbMsg := &model.Message{
				MsgId:      msg.MsgId,
				SeqId:      msg.SeqId,
				MsgType:    msg.MsgType,
				Content:    msg.Content,
				FromType:   msg.FromType,
				FromId:     msg.FromId,
				FromName:   msg.FromName,
				FromDevice: msg.FromDevice,
				ToType:     msg.ToType,
				ToId:       msg.ToId,
				ToName:     msg.ToName,
				ToDevice:   msg.ToDevice,
				SendTime:   msg.SendTime,
			}

			err = mainDb.SaveMessage(dbMsg)
			if err != nil {
				global.Logger.Error("Error inserting into Tidb", zap.Error(err))
				m.Requeue(-1)
				return
			}

			global.Logger.Info("Store message", zap.String("msgId", msg.MsgId), zap.String("fromId", msg.FromId), zap.Int64("seqId", msg.SeqId))

			m.Finish()
		}
	}(m))

	return nil
}
