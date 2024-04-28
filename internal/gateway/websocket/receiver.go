package websocket

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/model"
	"eim/internal/mq"
	"eim/internal/protocol"
	"eim/pkg/json"
	"eim/pkg/log"
	"eim/proto/pb"
)

func (its *Server) receiverHandler(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	start := time.Now()
	cmd, frame := protocol.WebsocketCodec.Decode(data)

	sess := conn.Session().(*session)

	switch cmd {
	case protocol.Message:
		{
			pbMsg := &pb.Message{}
			err := proto.Unmarshal(frame, pbMsg)
			if err != nil {
				atomic.AddInt64(&its.invalidMsgTotal, 1)
				log.Error("Illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			id := ""
			if pbMsg.ToType == model.ToUser {
				id = pbMsg.FromId
			} else {
				id = pbMsg.ToId
			}

			pbMsg.SeqId, err = its.seqRpc.Number(id)
			if err != nil {
				log.Error("Error getting seq id: %v，%v", zap.String("id", id), zap.Error(err))
				return
			}

			pbMsg.SendTime = time.Now().UnixNano()

			its.workerPool.Go(func(sess *session, pbMsg *pb.Message) func() {
				return func() {
					body, err := json.Marshal(pbMsg)
					if err != nil {
						log.Error("Error serializing message", zap.Error(err))
						return
					}

					if topic, exist := mq.MessageTopics[strconv.FormatInt(pbMsg.ToType, 10)]; exist {
						err = its.producer.Publish(string(topic), body)
						if err != nil {
							log.Error("Error publishing message", zap.Error(err))
							return
						}
						sess.send(protocol.Ack, []byte(pbMsg.MsgId))
					} else {
						atomic.AddInt64(&its.invalidMsgTotal, 1)
						log.Error("Illegal message", zap.ByteString("body", body))
					}
				}
			}(sess, pbMsg))

			atomic.AddInt64(&its.receivedTotal, 1)

			log.Debug("Time consuming to process messages", zap.Duration("duration", time.Since(start)))
		}
	}
}
