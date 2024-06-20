package websocket

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/gateway/protocol"
	"eim/internal/model"
	"eim/util/log"
)

type SendHandler struct {
	Server *Server
}

func (its *SendHandler) HandleMessage(m *nats.Msg) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	if m.Data == nil || len(m.Data) == 0 {
		_ = m.Ack()
		return fmt.Errorf("message data is nil")
	}

	msg := &model.Message{}
	err := proto.Unmarshal(m.Data, msg)
	if err != nil {
		_ = m.Ack()
		return fmt.Errorf("unmarshal message -> %w", err)
	}

	_ = m.Ack()

	var allSession []*session

	sessions := its.Server.sessionManager.Get(msg.UserId)
	for _, sess := range sessions {
		allSession = append(allSession, sess)
	}

	for _, s := range allSession {
		s.send(protocol.Message, m.Data)
	}

	atomic.AddInt64(&its.Server.sendMsgTotal, 1)

	return nil
}
