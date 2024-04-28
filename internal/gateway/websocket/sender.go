package websocket

import (
	"sync/atomic"

	"github.com/nsqio/go-nsq"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/protocol"
	"eim/pkg/log"
)

type SendHandler struct {
	Server   *Server
	TaskPool *ants.Pool
}

func (its *SendHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		atomic.AddInt64(&its.Server.invalidMsgTotal, 1)
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

			var allSession []*session

			sessions := its.Server.sessionManager.Get(msg.FromId)
			for _, session := range sessions {
				allSession = append(allSession, session)
			}

			for _, s := range allSession {
				s.send(protocol.Message, m.Body)
			}

			atomic.AddInt64(&its.Server.sendTotal, 1)

			m.Finish()
		}
	}(m))
	if err != nil {
		return err
	}

	return nil
}
