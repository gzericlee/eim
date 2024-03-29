package websocket

import (
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/protocol"
	"eim/internal/types"
	"eim/pkg/log"
)

type session struct {
	device *types.Device
	conn   *websocket.Conn
}

func (its *session) send(cmd int, body []byte) {
	err := its.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(cmd, body))
	if err != nil {
		log.Error("Error sending message", zap.String("userId", its.device.UserId), zap.String("deviceId", its.device.DeviceId), zap.Error(err))
		_ = its.conn.Close()
		return
	}
	log.Debug("Sent message successful", zap.String("userId", its.device.UserId), zap.String("deviceId", its.device.DeviceId))
}

type manager struct {
	sync.Map
}

func (its *manager) Add(userId string, sess *session) {
	var sessions = map[string]*session{}
	if value, exist := its.Load(userId); exist {
		sessions = value.(map[string]*session)
	}
	sessions[sess.device.DeviceId] = sess
	its.Store(userId, sessions)
}

func (its *manager) Get(userId string) map[string]*session {
	if value, exist := its.Load(userId); exist {
		return value.(map[string]*session)
	}
	return nil
}

func (its *manager) Remove(userId, deviceId string) {
	if value, exist := its.Load(userId); exist {
		sessions := value.(map[string]*session)
		delete(sessions, deviceId)
	}
}
