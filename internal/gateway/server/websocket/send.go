package websocket

import (
	"encoding/json"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/protocol"
	"eim/internal/gateway/session"
	"eim/util/log"
)

func (its *Server) Send(sess *session.Session, cmd int, body []byte) {
	its.send(sess, cmd, body)
}

func (its *Server) sendOfflineMessage(sess *session.Session) {
	if sess == nil {
		return
	}

	device := sess.GetDevice()

	messages, err := its.storageRpc.GetOfflineMessages(device.UserId, device.DeviceId)
	if err != nil {
		log.Error("Error get offline messages", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.Error(err))
		return
	}
	body, err := json.Marshal(messages)
	if err != nil {
		log.Error("Error marshal offline messages", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.Error(err))
		return
	}
	its.send(sess, protocol.Messages, body)
}

func (its *Server) send(sess *session.Session, cmd int, body []byte) {
	if sess == nil {
		return
	}

	device := sess.GetDevice()
	conn := sess.GetConn().(*websocket.Conn)

	err := conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(cmd, body))
	if err != nil {
		log.Error("Error send message", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.Error(err))
		return
	}

	log.Debug("Send message successfully", zap.Int("cmd", cmd), zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId))
}
