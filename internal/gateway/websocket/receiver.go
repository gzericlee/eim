package websocket

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/nsq/producer"
	"eim/internal/protocol"
	"eim/internal/types"
	"eim/pkg/json"
	"eim/pkg/log"
	"eim/proto/pb"
)

func receiverHandler(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(gatewaySvr.keepaliveTime))

	start := time.Now()
	cmd, frame := protocol.WebsocketCodec.Decode(data)

	sess := conn.Session().(*session)

	switch cmd {
	case protocol.Message:
		{
			pbMsg := &pb.Message{}
			err := proto.Unmarshal(frame, pbMsg)
			if err != nil {
				log.Error("Illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			id := ""
			if pbMsg.ToType == types.ToUser {
				id = pbMsg.FromId
			} else {
				id = pbMsg.ToId
			}

			pbMsg.SeqId, err = seqRpc.Number(id)
			if err != nil {
				log.Error("Error getting seq id: %v，%v", zap.String("id", id), zap.Error(err))
				return
			}

			pbMsg.SendTime = time.Now().UnixNano()

			gatewaySvr.workerPool.Go(func(sess *session, pbMsg *pb.Message) func() {
				return func() {
					body, err := json.Marshal(pbMsg)
					if err != nil {
						log.Error("Error serializing message", zap.Error(err))
						return
					}

					topic := ""
					switch pbMsg.ToType {
					case types.ToUser:
						topic = types.MessageUserDispatchTopic
					case types.ToGroup:
						topic = types.MessageGroupDispatchTopic
					}

					err = producer.PublishAsync(topic, body)
					if err != nil {
						log.Error("Error publishing message", zap.Error(err))
						return
					}

					sess.send(protocol.Ack, []byte(pbMsg.MsgId))
				}
			}(sess, pbMsg))

			gatewaySvr.receivedTotal.Add(1)

			log.Debug("Time consuming to process messages", zap.Duration("duration", time.Since(start)))
		}
	}
}
