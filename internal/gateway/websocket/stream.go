package websocket

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/nsq/producer"
	"eim/internal/protocol"
	"eim/model"
	"eim/proto/pb"
)

func streamHandler(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
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
				global.Logger.Error("Illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			uId := ""
			if pbMsg.ToType == model.ToUser {
				uId = pbMsg.FromId
			} else {
				uId = pbMsg.ToId
			}
			pbMsg.SeqId, err = seqRpc.Id(uId)
			if err != nil {
				global.Logger.Error("Error getting seq id: %v，%v", zap.String("userId", pbMsg.FromId), zap.Error(err))
				return
			}

			pbMsg.SendTime = time.Now().UnixNano()

			gatewaySvr.workerPool.Go(func(sess *session, pbMsg *pb.Message) func() {
				return func() {
					frame, err := proto.Marshal(pbMsg)
					if err != nil {
						global.Logger.Error("Error serializing message", zap.Error(err))
						return
					}
					err = producer.PublishAsync(model.MessageDispatchTopic, frame)
					if err != nil {
						global.Logger.Error("Error publishing message", zap.Error(err))
						return
					}
					sess.send(protocol.Ack, []byte(pbMsg.MsgId))
				}
			}(sess, pbMsg))

			gatewaySvr.receivedTotal.Add(1)

			global.Logger.Debug("Time consuming to process messages", zap.Duration("duration", time.Since(start)))
		}
	}
}
