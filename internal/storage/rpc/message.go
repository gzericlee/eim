package rpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eim/internal/model"
	"eim/util/log"

	"eim/internal/database"
)

type Message struct {
	Database database.IDatabase
}

type MessageRequest struct {
	Message *model.Message
}

type MessageReply struct {
}

func (its *Message) Save(ctx context.Context, req *MessageRequest, reply *MessageReply) error {
	err := its.Database.SaveMessage(req.Message)
	if err != nil {
		return fmt.Errorf("save db message -> %w", err)
	}

	log.Debug("Store message", zap.String("msgId", req.Message.MsgId), zap.Int64("seqId", req.Message.SeqId))

	return nil
}
