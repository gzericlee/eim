package websocket

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eim/internal/gateway/session"
	"eim/internal/model"
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

	device.OfflineAt = timestamppb.Now()
	device.State = model.OfflineState

	err = its.storageRpc.SaveDevice(device)
	if err != nil {
		log.Error("Error save device", zap.Error(err))
		return
	}
}
