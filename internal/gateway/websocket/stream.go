package websocket

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"

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
				global.Logger.Errorf("Illegal message: %v，%v", frame, err)
				return
			}

			uId := ""
			if pbMsg.ToType == model.ToUser {
				uId = pbMsg.FromId
			} else {
				uId = pbMsg.ToId
			}
			pbMsg.SeqId, err = seqSvr.ID(uId)
			if err != nil {
				global.Logger.Errorf("Error geting seq id: %v，%v", pbMsg.FromId, err)
				return
			}

			pbMsg.SendTime = time.Now().UnixNano()

			gatewaySvr.workerPool.Go(func(sess *session, pbMsg *pb.Message) func() {
				return func() {
					frame, err := proto.Marshal(pbMsg)
					if err != nil {
						global.Logger.Warnf("Error serializing message: %v", err)
						return
					}
					err = producer.PublishAsync(model.MessageDispatchTopic, frame)
					if err != nil {
						global.Logger.Warnf("Error publishing message: %v", err)
						return
					}
					sess.send(protocol.Ack, []byte(pbMsg.MsgId))
				}
			}(sess, pbMsg))

			gatewaySvr.receivedTotal.Add(1)

			global.Logger.Debugf("Time consuming to process messages: %v", time.Since(start))
		}
	}
}
