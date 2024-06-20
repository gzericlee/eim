package websocket

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/protocol"
	"eim/internal/model"
	"eim/internal/mq"
	"eim/util/log"
)

func (its *Server) receiverHandler(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	cmd, frame := protocol.WebsocketCodec.Decode(data)

	sess := conn.Session().(*session)

	switch cmd {
	case protocol.Ack:
		{
			msgId := string(frame)
			err := its.storageRpc.RemoveOfflineMessageIds([]interface{}{msgId}, sess.user.UserId, sess.device.DeviceId)
			if err != nil {
				atomic.AddInt64(&its.errorTotal, 1)
				log.Error("Error remove offline message id", zap.Error(err))
				return
			}
			atomic.AddInt64(&its.ackTotal, 1)
		}
	case protocol.Message:
		{
			msg := &model.Message{}
			err := proto.Unmarshal(frame, msg)
			if err != nil {
				atomic.AddInt64(&its.invalidMsgTotal, 1)
				log.Error("Error illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			bizId := ""
			if msg.ToType == model.ToUser {
				bizId = msg.FromId
			} else {
				bizId = msg.ToId
			}

			msg.SeqId, err = its.seqRpc.IncrementId(bizId)
			if err != nil {
				log.Error("Error getting seq bizId: %v，%v", zap.String("bizId", bizId), zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}

			msg.SendTime = time.Now().UnixNano()

			body, err := proto.Marshal(msg)
			if err != nil {
				log.Error("Error marshal message", zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}

			switch msg.ToType {
			case model.ToUser:
				err = its.producer.Publish(mq.UserMessageSubject, body)
			case model.ToGroup:
				err = its.producer.Publish(mq.GroupMessageSubject, body)
			}
			if err != nil {
				log.Error("Error publish message", zap.Error(err))
				atomic.AddInt64(&its.errorTotal, 1)
				return
			}

			sess.send(protocol.Ack, []byte(strconv.FormatInt(msg.MsgId, 10)))

			atomic.AddInt64(&its.receivedMsgTotal, 1)
		}
	}
}
