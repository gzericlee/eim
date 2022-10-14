package websocket

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/internal/pool"
	"eim/internal/protocol"
	"eim/internal/types"
	"eim/pkg/log"
)

var workerCount = make(chan struct{}, 100000)

type SendHandler struct {
}

func (its *SendHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		gatewaySvr.invalidMsgTotal.Add(1)
		return nil
	}

	pool.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			msg := &types.Message{}
			err := msg.Deserialize(m.Body)
			if err != nil {
				log.Error("Error deserializing message", zap.Error(err))
				m.Finish()
				return
			}

			var allSession []*session

			sessions := gatewaySvr.sessionManager.Get(msg.FromId)
			for _, session := range sessions {
				allSession = append(allSession, session)
			}

			for _, s := range allSession {
				//if session.device.DeviceId == msg.FromDevice {
				//	continue
				//}
				workerCount <- struct{}{}
				go func(s *session, body []byte) {
					s.send(protocol.Message, body)
					<-workerCount
				}(s, m.Body)
			}

			gatewaySvr.sentTotal.Add(1)

			m.Finish()
		}
	}(m))

	return nil
}
