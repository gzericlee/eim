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

	"eim/global"
	"eim/internal/nsq/producer"
	"eim/internal/seq"
	"eim/model"
)

var gatewaySvr *server
var seqSvr *seq.RpcClient

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

	u := websocket.NewUpgrader()

	u.OnMessage(streamHandler)

	u.SetPingHandler(func(conn *websocket.Conn, s string) {
		_ = conn.SetReadDeadline(time.Now().Add(gatewaySvr.keepaliveTime))
		_ = conn.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
		gatewaySvr.heartbeatTotal.Add(1)
	})

	u.OnClose(func(conn *websocket.Conn, err error) {
		sess := conn.Session().(*session)
		if sess == nil {
			return
		}

		defer func() {
			gatewaySvr.clientTotal.Add(-1)
			gatewaySvr.sessionManager.Remove(sess.device.DeviceId)
			global.Logger.Debugf("Device disconnected: %s %v", sess.device.DeviceId, err)
		}()

		if !sess.verified {
			return
		}

		now := time.Now().Local()
		sess.device.OfflineAt = &now
		sess.device.State = model.OfflineState

		gatewaySvr.workerPool.Go(func(device *model.Device) func() {
			return func() {
				body, _ := device.Serialize()
				err = producer.PublishAsync(model.DeviceStoreTopic, body)
				if err != nil {
					global.Logger.Warnf("Error publishing message: %v", err)
					return
				}
			}
		}(sess.device))
	})

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		global.Logger.Errorf("Error upgrading websocket protocol: %s", err.Error())
		return
	}
	wsConn := conn.(*websocket.Conn)
	_ = wsConn.SetReadDeadline(time.Now().Add(gatewaySvr.keepaliveTime))

	//TODO 校验身份
	sess := &session{device: &model.Device{}}
	now := time.Now().Local()
	sess.device.OnlineAt = &now
	sess.device.DeviceId = r.Header.Get("Deviceid")
	sess.device.UserId = r.Header.Get("Userid")
	sess.device.DeviceVersion = r.Header.Get("Deviceversion")
	sess.device.DeviceType = r.Header.Get("Devicetype")
	sess.device.State = model.OnlineState
	sess.device.GatewayIp = global.SystemConfig.LocalIp
	sess.verified = true
	sess.conn = wsConn

	wsConn.SetSession(sess)

	sessions := gatewaySvr.sessionManager.GetByUserId(sess.device.UserId)
	sessions = append(sessions, sess)

	gatewaySvr.sessionManager.Save(sess.device.UserId, sessions)

	gatewaySvr.workerPool.Go(func(device *model.Device) func() {
		return func() {
			body, _ := device.Serialize()
			err = producer.PublishAsync(model.DeviceStoreTopic, body)
			if err != nil {
				global.Logger.Warnf("Error publishing message: %v", err)
				_ = wsConn.Close()
				return
			}
		}
	}(sess.device))

	gatewaySvr.clientTotal.Add(1)

	global.Logger.Debugf("Device login successful: %s，%s，%s", sess.device.UserId, sess.device.DeviceId, sess.device.DeviceVersion)
}

func InitGatewayServer(ip string, ports []string) error {
	logging.SetLevel(logging.LevelNone)

	var err error
	seqSvr, err = seq.NewRpcClient(global.SystemConfig.SeqSvr.Endpoint)
	if err != nil {
		return err
	}

	pool := taskpool.NewFixedPool(runtime.NumCPU()*5, 1024)
	//pool, err := ants.NewPool(runtime.NumCPU()*1000, ants.WithNonblocking(true), ants.WithPreAlloc(true), ants.WithPanicHandler(func(i interface{}) {
	//	global.Logger.Errorf("Panic with worker pool: %v", i)
	//}))

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

	global.Logger.Infof("Listening Websocket from %v : %v", ip, ports)

	return nil
}
