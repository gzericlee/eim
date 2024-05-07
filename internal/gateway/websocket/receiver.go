package websocket

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/protocol"
	"eim/internal/model"
	"eim/internal/mq"
	"eim/pkg/log"
)

func (its *Server) receiverHandler(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	start := time.Now()
	cmd, frame := protocol.WebsocketCodec.Decode(data)

	sess := conn.Session().(*session)

	switch cmd {
	case protocol.Message:
		{
			msg := &model.Message{}
			err := proto.Unmarshal(frame, msg)
			if err != nil {
				atomic.AddInt64(&its.invalidMsgTotal, 1)
				log.Error("Illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			id := ""
			if msg.ToType == model.ToUser {
				id = msg.FromId
			} else {
				id = msg.ToId
			}

			msg.SeqId, err = its.seqRpc.IncrementId(id)
			if err != nil {
				log.Error("Error getting seq id: %v，%v", zap.String("id", id), zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}

			msgId, err := its.seqRpc.SnowflakeId()
			if err != nil {
				log.Error("Error getting snowflake id", zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}
			msg.MsgId = msgId
			msg.SendTime = time.Now().UnixNano()

			err = its.workerPool.Submit(func(sess *session, pbMsg *model.Message) func() {
				return func() {
					body, err := proto.Marshal(pbMsg)
					if err != nil {
						log.Error("Error serializing message", zap.Error(err))
						atomic.AddInt64(&its.errorTotal, 1)
						return
					}

					err = its.producer.Publish(mq.MessageDispatchSubject, body)
					if err != nil {
						log.Error("Error publishing message", zap.Error(err))
						atomic.AddInt64(&its.errorTotal, 1)
						return
					}
					idBody := make([]byte, binary.MaxVarintLen64)
					binary.PutVarint(idBody, pbMsg.MsgId)
					sess.send(protocol.Ack, idBody)
				}
			}(sess, msg))

			if err != nil {
				log.Error("Error submitting task", zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}

			atomic.AddInt64(&its.receivedMsgTotal, 1)

			log.Debug("Time consuming to process messages", zap.Duration("duration", time.Since(start)))
		}
	}
}
