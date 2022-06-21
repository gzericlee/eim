package websocket

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/lesismal/nbio/taskpool"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/auth"
	"eim/internal/nsq/producer"
	"eim/internal/seq"
	"eim/model"
)

var gatewaySvr *server
var seqRpc *seq.RpcClient
var authRpc *auth.RpcClient

type server struct {
	ports           []string
	sessionManager  *manager
	receivedTotal   *atomic.Int64
	sentTotal       *atomic.Int64
	invalidMsgTotal *atomic.Int64
	heartbeatTotal  *atomic.Int64
	clientTotal     *atomic.Int64
	workerPool      *taskpool.FixedPool
	keepaliveTime   time.Duration
	http            *nbhttp.Server
}

func (its *server) connHandler(w http.ResponseWriter, r *http.Request) {
	var isWs bool
	if strings.ToUpper(r.Header.Get("Connection")) == "UPGRADE" && strings.ToUpper(r.Header.Get("Upgrade")) == "WEBSOCKET" {
		isWs = true
	}
	if !isWs {
		_, _ = w.Write([]byte("Only websocket connections are supported"))
		return
	}

	user, err := authRpc.CheckToken(r.Header.Get("Token"))
	if err != nil {
		global.Logger.Error("Error checking token", zap.Error(err))
		return
	}

	u := websocket.NewUpgrader()

	u.OnMessage(receiverHandler)

	u.SetPingHandler(func(conn *websocket.Conn, s string) {
		_ = conn.SetReadDeadline(time.Now().Add(gatewaySvr.keepaliveTime))
		_ = conn.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
		gatewaySvr.heartbeatTotal.Add(1)
	})

	u.OnClose(func(conn *websocket.Conn, err error) {
		session := conn.Session().(*session)
		if session == nil {
			return
		}

		defer func() {
			gatewaySvr.clientTotal.Add(-1)
			gatewaySvr.sessionManager.Remove(session.device.UserId, session.device.DeviceId)
			global.Logger.Debug("Device disconnected", zap.String("deviceId", session.device.DeviceId), zap.Error(err))
		}()

		now := time.Now().Local()
		session.device.OfflineAt = &now
		session.device.State = model.OfflineState

		gatewaySvr.workerPool.Go(func(device *model.Device) func() {
			return func() {
				body, _ := device.Serialize()
				err = producer.PublishAsync(model.DeviceStoreTopic, body)
				if err != nil {
					global.Logger.Error("Error publishing device", zap.Error(err))
				}
			}
		}(session.device))
	})

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		global.Logger.Error("Error upgrading websocket protocol", zap.Error(err))
		return
	}

	wsConn := conn.(*websocket.Conn)
	_ = wsConn.SetReadDeadline(time.Now().Add(gatewaySvr.keepaliveTime))

	session := &session{device: &model.Device{}}
	now := time.Now().Local()

	//TODO 为了方便模拟，这里直接取Header的UserId，实际应该取Auth服务返回的User
	session.device.UserId = r.Header.Get("Userid")
	if session.device.UserId == "" {
		session.device.UserId = user.UserId
	}

	session.device.OnlineAt = &now
	session.device.DeviceId = r.Header.Get("Deviceid")
	session.device.DeviceVersion = r.Header.Get("Deviceversion")
	session.device.DeviceType = r.Header.Get("Devicetype")
	session.device.State = model.OnlineState
	session.device.GatewayIp = global.SystemConfig.LocalIp
	session.conn = wsConn

	wsConn.SetSession(session)

	gatewaySvr.sessionManager.Add(session.device.UserId, session)

	gatewaySvr.workerPool.Go(func(device *model.Device) func() {
		return func() {
			body, _ := device.Serialize()
			err = producer.PublishAsync(model.DeviceStoreTopic, body)
			if err != nil {
				global.Logger.Error("Error publishing device", zap.Error(err))
			}
			gatewaySvr.clientTotal.Add(1)
			global.Logger.Debug("Device login successful", zap.String("userId", device.UserId), zap.String("deviceId", device.DeviceId), zap.String("version", device.DeviceVersion))
		}
	}(session.device))
}

func InitWebsocketServer(ip string, ports []string) error {
	logging.SetLevel(logging.LevelNone)

	var err error
	seqRpc, err = seq.NewRpcClient(global.SystemConfig.Etcd.Endpoints.Value())
	if err != nil {
		return err
	}

	authRpc, err = auth.NewRpcClient(global.SystemConfig.Etcd.Endpoints.Value())
	if err != nil {
		return err
	}

	pool := taskpool.NewFixedPool(runtime.NumCPU(), 1024)

	gatewaySvr = &server{
		ports:           ports,
		workerPool:      pool,
		sentTotal:       &atomic.Int64{},
		receivedTotal:   &atomic.Int64{},
		heartbeatTotal:  &atomic.Int64{},
		invalidMsgTotal: &atomic.Int64{},
		clientTotal:     &atomic.Int64{},
		sessionManager:  new(manager),
		keepaliveTime:   time.Minute * 5,
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", gatewaySvr.connHandler)

	var addrs []string
	for _, port := range ports {
		addrs = append(addrs, fmt.Sprintf("%v:%v", ip, port))
	}

	gatewaySvr.http = nbhttp.NewServer(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   addrs,
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
		KeepaliveTime:           gatewaySvr.keepaliveTime,
	})

	err = gatewaySvr.http.Start()
	if err != nil {
		return err
	}

	printWebSocketServiceDetail()

	global.Logger.Info("Listening websocket connect", zap.String("ip", ip), zap.Strings("ports", ports))

	return nil
}
