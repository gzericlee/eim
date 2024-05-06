package websocket

import (
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/protocol"
	"eim/pkg/log"
)

type SendHandler struct {
	Server *Server
}

func (its *SendHandler) HandleMessage(data []byte) error {
	if data == nil || len(data) == 0 {
		atomic.AddInt64(&its.Server.invalidMsgTotal, 1)
		return nil
	}

	msg := &model.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("Error deserializing message", zap.Error(err))
		return nil
	}

	var allSession []*session

	sessions := its.Server.sessionManager.Get(msg.UserId)
	for _, session := range sessions {
		allSession = append(allSession, session)
	}

	for _, s := range allSession {
		s.send(protocol.Message, data)
	}

	atomic.AddInt64(&its.Server.sendMsgTotal, 1)

	return nil
}
