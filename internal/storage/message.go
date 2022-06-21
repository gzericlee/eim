package storage

import (
	"context"

	"go.uber.org/zap"

	"eim/global"
	"eim/model"
)

type Message struct {
}

type MessageRequest struct {
	Message *model.Message
}

type MessageReply struct {
}

func (its *Message) Save(ctx context.Context, req *MessageRequest, reply *MessageReply) error {
	err := mainDb.SaveMessage(req.Message)
	if err != nil {
		global.Logger.Error("Error inserting into Tidb", zap.Error(err))
		return err
	}

	global.Logger.Info("Store message", zap.String("msgId", req.Message.MsgId), zap.String("fromId", req.Message.FromId), zap.Int64("seqId", req.Message.SeqId))

	return nil
}
