package session

import (
	"github.com/lesismal/nbio/nbhttp/websocket"

	"github.com/gzericlee/eim/internal/model"
)

type Session struct {
	user   *model.Biz
	device *model.Device
	conn   interface{}
}

func NewSession(user *model.Biz, device *model.Device, conn *websocket.Conn) *Session {
	return &Session{
		user:   user,
		device: device,
		conn:   conn,
	}
}

func (its *Session) GetUser() *model.Biz {
	return its.user
}

func (its *Session) GetDevice() *model.Device {
	return its.device
}

func (its *Session) GetConn() interface{} {
	return its.conn
}
