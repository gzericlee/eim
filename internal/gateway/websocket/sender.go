package websocket

import (
	"fmt"
	"sync/atomic"

	"github.com/golang/protobuf/proto"

	"eim/internal/gateway/protocol"
	"eim/internal/model"
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
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	var allSession []*session

	sessions := its.Server.sessionManager.Get(msg.UserId)
	for _, sess := range sessions {
		allSession = append(allSession, sess)
	}

	for _, s := range allSession {
		s.send(protocol.Message, data)
	}

	atomic.AddInt64(&its.Server.sendMsgTotal, 1)

	return nil
}
