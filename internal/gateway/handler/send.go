package handler

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"

	"eim/internal/gateway/protocol"
	"eim/internal/gateway/server"
	"eim/internal/gateway/session"
	"eim/internal/model"
)

type SendHandler struct {
	Servers []server.IServer
}

func (its *SendHandler) Process(m *nats.Msg) error {
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

	var allSession []*session.Session

	for _, server := range its.Servers {
		sessions := server.GetSessionManager().Get(msg.UserId)
		for _, sess := range sessions {
			allSession = append(allSession, sess)
		}

		for _, sess := range allSession {
			server.Send(sess, protocol.Message, m.Data)
		}

		server.IncrSendMsgTotal(1)
	}

	return m.Ack()
}
