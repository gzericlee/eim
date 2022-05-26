package websocket

import (
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/protocol"
	"eim/model"
	"eim/proto/pb"
)

type SendHandler struct {
}

func (its *SendHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		gatewaySvr.invalidMsgTotal.Add(1)
		return nil
	}

	global.SystemPool.Go(func(m *nsq.Message) func() {
		return func() {
			msg := &pb.Message{}
			err := proto.Unmarshal(m.Body, msg)
			if err != nil {
				global.Logger.Error("Error deserializing message", zap.Error(err))
				m.Finish()
				return
			}

			switch msg.ToType {
			case model.ToUser:
				{
					fromSess := gatewaySvr.sessionManager.GetByUserId(msg.FromId)
					toSess := gatewaySvr.sessionManager.GetByUserId(msg.ToId)

					allSess := append([]*session{}, fromSess...)
					allSess = append([]*session{}, toSess...)

					_ = allSess

					for _, sess := range allSess {
						//if session.device.DeviceId == msg.FromDevice {
						//	continue
						//}
						sess.send(protocol.Message, m.Body)
					}

					gatewaySvr.sentTotal.Add(1)
				}
			}

			m.Finish()
		}
	}(m))

	return nil
}
