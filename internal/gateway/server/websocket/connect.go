package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/gateway/protocol"
	"github.com/gzericlee/eim/internal/gateway/session"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/model/consts"
	"github.com/gzericlee/eim/pkg/log"
)

func (its *Server) connect(w http.ResponseWriter, r *http.Request) {
	var isWs bool
	if strings.ToUpper(r.Header.Get("Connection")) == "UPGRADE" && strings.ToUpper(r.Header.Get("Upgrade")) == "WEBSOCKET" {
		isWs = true
	}
	if !isWs {
		_, _ = w.Write([]byte("Only websocket connections are supported"))
		return
	}

	ws := websocket.NewUpgrader()

	ws.OnMessage(its.receive)
	ws.OnClose(its.close)
	ws.SetPingHandler(its.heartbeat)

	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error upgrade websocket protocol", zap.Error(err))
		return
	}

	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Replace(token, "Basic ", "", 1)

	user, err := its.authRpc.CheckToken(token)
	if err != nil {
		_ = conn.Close()
		log.Error("Error check auth token", zap.Error(err))
		return
	}

	gatewayAddr := fmt.Sprintf("%s:%d", its.ip, its.port)
	deviceId := r.Header.Get("DeviceId")
	deviceType := r.Header.Get("DeviceType")
	deviceVersion := r.Header.Get("DeviceVersion")

	device := &model.Device{
		DeviceId:      deviceId,
		UserId:        user.BizId,
		TenantId:      user.TenantId,
		DeviceType:    deviceType,
		DeviceVersion: deviceVersion,
		GatewayAddr:   gatewayAddr,
		OnlineAt:      time.Now().Unix(),
		State:         consts.StatusOnline,
	}
	err = its.deviceRpc.SaveDevice(device)
	if err != nil {
		_ = conn.Close()
		log.Error("Error save device", zap.Error(err))
		return
	}

	sess := session.NewSession(user, device, conn)
	conn.SetSession(sess)

	//预加载
	its.preload(sess)

	//离线消息
	its.sendOfflineMessage(sess)

	its.sessionManager.Add(user.BizId, sess)

	its.send(sess, protocol.Connected, nil)

	its.IncrClientTotal(1)

	//log.Debug("Device connected successfully", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.String("version", device.DeviceVersion))
}
