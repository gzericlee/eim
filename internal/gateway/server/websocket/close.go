package websocket

import (
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/session"
	"eim/internal/model/consts"
	"eim/pkg/log"
)

func (its *Server) close(conn *websocket.Conn, err error) {
	sess := conn.Session().(*session.Session)
	if sess == nil {
		return
	}

	device := sess.GetDevice()

	defer func() {
		its.IncrClientTotal(-1)
		its.sessionManager.Remove(device.UserId, device.DeviceId)
		//log.Warn("Device disconnected", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.String("version", device.DeviceVersion))
	}()

	device.OfflineAt = time.Now().Unix()
	device.State = consts.StatusOffline

	err = its.storageRpc.UpdateDevice(device)
	if err != nil {
		log.Error("Error update device", zap.Error(err))
		return
	}
}
