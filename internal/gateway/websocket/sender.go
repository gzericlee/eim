package websocket

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/protocol"
	"eim/internal/redis"
	"eim/model"
)

var workerCount = make(chan struct{}, 100000)

type SendHandler struct {
}

func (its *SendHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		gatewaySvr.invalidMsgTotal.Add(1)
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			msg := &model.Message{}
			err := msg.Deserialize(m.Body)
			if err != nil {
				global.Logger.Error("Error deserializing message", zap.Error(err))
				m.Finish()
				return
			}

			var allSession []*session

			sessions := gatewaySvr.sessionManager.Get(msg.FromId)
			for _, session := range sessions {
				allSession = append(allSession, session)
			}

			switch msg.ToType {
			case model.ToUser:
				{
					sessions := gatewaySvr.sessionManager.Get(msg.ToId)
					for _, session := range sessions {
						allSession = append(allSession, session)
					}
				}
			case model.ToGroup:
				{
					members, err := redis.GetGroupMembers(msg.ToId)
					if err != nil {
						global.Logger.Error("Error getting Group members", zap.String("groupId", msg.ToId), zap.Error(err))
						return
					}
					for _, userId := range members {
						sessions := gatewaySvr.sessionManager.Get(userId)
						for _, session := range sessions {
							allSession = append(allSession, session)
						}
					}
				}
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
