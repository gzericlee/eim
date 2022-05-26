package websocket

import (
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/protocol"
	"eim/model"
)

type session struct {
	device   *model.Device
	conn     *websocket.Conn
	verified bool
}

func (its *session) send(cmd int, body []byte) {
	err := its.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(cmd, body))
	if err != nil {
		global.Logger.Error("Error sending message", zap.String("userId", its.device.UserId), zap.String("deviceId", its.device.DeviceId), zap.Error(err))
		return
	}
	global.Logger.Debug("Sent message successful", zap.String("userId", its.device.UserId), zap.String("deviceId", its.device.DeviceId))
}

type manager struct {
	sync.Map
}

func (its *manager) Save(userId string, sessions []*session) {
	its.Store(userId, sessions)
}

func (its *manager) GetByUserId(userId string) []*session {
	if value, exist := its.Load(userId); exist {
		return value.([]*session)
	}
	return nil
}

func (its *manager) Remove(key string) {
	its.Delete(key)
}
