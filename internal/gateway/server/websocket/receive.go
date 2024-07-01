package websocket

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/protocol"
	"eim/internal/gateway/session"
	"eim/internal/model"
	"eim/internal/mq"
	"eim/util/log"
)

func (its *Server) receive(conn *websocket.Conn, _ websocket.MessageType, data []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	cmd, frame := protocol.WebsocketCodec.Decode(data)

	sess := conn.Session().(*session.Session)
	if sess == nil {
		return
	}

	user := sess.GetUser()
	device := sess.GetDevice()

	switch cmd {
	case protocol.Ack:
		{
			msgId := string(frame)
			err := its.storageRpc.RemoveOfflineMessages([]string{msgId}, user.BizId, device.DeviceId)
			if err != nil {
				its.IncrErrorTotal(1)
				log.Error("Error remove offline message id", zap.Error(err))
				return
			}
			its.IncrAckTotal(1)
		}
	case protocol.Message:
		{
			msg := &model.Message{}
			err := proto.Unmarshal(frame, msg)
			if err != nil {
				its.IncrInvalidMsgTotal(1)
				log.Error("Error illegal message", zap.ByteString("body", frame), zap.Error(err))
				return
			}

			msg.SeqId, err = its.seqRpc.IncrId(user.BizId, user.TenantId)
			if err != nil {
				log.Error("Error getting seq", zap.String("bizId", user.BizId), zap.String("tenantId", user.TenantId), zap.Error(err))
				its.IncrErrorTotal(1)
				return
			}

			msg.SendTime = time.Now().UnixNano()

			body, err := proto.Marshal(msg)
			if err != nil {
				log.Error("Error marshal message", zap.Error(err))
				its.IncrErrorTotal(1)
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
				its.IncrErrorTotal(1)
				return
			}

			its.send(sess, protocol.Ack, []byte(strconv.FormatInt(msg.MsgId, 10)))

			its.IncrReceivedMsgTotal(1)
		}
	}
}
