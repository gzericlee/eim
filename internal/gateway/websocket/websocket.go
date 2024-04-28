package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/lesismal/nbio/taskpool"
	"go.uber.org/zap"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/config"
	"eim/internal/model"
	"eim/internal/mq"
	"eim/internal/redis"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
)

type Server struct {
	ip             string
	ports          []string
	keepaliveTime  time.Duration
	sessionManager *manager
	workerPool     *taskpool.TaskPool
	http           *nbhttp.Server
	seqRpc         *seqrpc.Client
	authRpc        *authrpc.Client
	storageRpc     *storagerpc.Client
	redisManager   *redis.Manager
	producer       mq.Producer

	receivedTotal   int64
	sendTotal       int64
	invalidMsgTotal int64
	heartbeatTotal  int64
	clientTotal     int64
}

func NewServer(ip string, ports []string, seqRpc *seqrpc.Client, authRpc *authrpc.Client, storageRpc *storagerpc.Client, redisManager *redis.Manager, producer mq.Producer) *Server {
	var address []string
	for _, port := range ports {
		address = append(address, fmt.Sprintf("%v:%v", ip, port))
	}

	keepaliveTime := time.Minute * 5

	server := &Server{
		ip:             ip,
		ports:          ports,
		keepaliveTime:  keepaliveTime,
		sessionManager: new(manager),
		workerPool:     taskpool.New(32, 1024),
		seqRpc:         seqRpc,
		authRpc:        authRpc,
		redisManager:   redisManager,
		storageRpc:     storageRpc,
		producer:       producer,
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", server.connectHandler)

	server.http = nbhttp.NewServer(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   address,
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
		KeepaliveTime:           keepaliveTime,
	})

	go server.printServiceDetail()

	return server
}

func (its *Server) Start() error {
	return its.http.Start()
}

func (its *Server) connectHandler(w http.ResponseWriter, r *http.Request) {
	var isWs bool
	if strings.ToUpper(r.Header.Get("Connection")) == "UPGRADE" && strings.ToUpper(r.Header.Get("Upgrade")) == "WEBSOCKET" {
		isWs = true
	}
	if !isWs {
		_, _ = w.Write([]byte("Only websocket connections are supported"))
		return
	}

	user, err := its.authRpc.CheckToken(r.Header.Get("Token"))
	//if err != nil {
	//	log.Error("Error checking token", zap.Error(err))
	//	return
	//}

	ws := websocket.NewUpgrader()

	ws.OnMessage(its.receiverHandler)

	ws.SetPingHandler(func(conn *websocket.Conn, s string) {
		_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))
		_ = conn.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
		atomic.AddInt64(&its.heartbeatTotal, 1)
	})

	ws.OnClose(func(conn *websocket.Conn, err error) {
		sess := conn.Session().(*session)
		if sess == nil {
			return
		}

		defer func() {
			atomic.AddInt64(&its.clientTotal, -1)
			its.sessionManager.Remove(sess.device.UserId, sess.device.DeviceId)
			log.Debug("Device disconnected", zap.String("deviceId", sess.device.DeviceId), zap.Error(err))
		}()

		now := time.Now().Local()
		sess.device.OfflineAt = &now
		sess.device.State = model.OfflineState

		err = its.storageRpc.SaveDevice(sess.device)
		if err != nil {
			log.Error("Error saving device", zap.Error(err))
			return
		}
	})

	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error upgrading websocket protocol", zap.Error(err))
		return
	}

	_ = conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))

	now := time.Now().Local()

	sess := &session{device: &model.Device{}}

	//TODO 为了方便模拟，这里直接取Header的UserId，实际应该取Auth服务返回的User
	sess.device.UserId = r.Header.Get("UserId")
	if sess.device.UserId == "" {
		sess.device.UserId = user.UserId
	}

	sess.device.OnlineAt = &now
	sess.device.DeviceId = r.Header.Get("DeviceId")
	sess.device.DeviceVersion = r.Header.Get("DeviceVersion")
	sess.device.DeviceType = r.Header.Get("DeviceType")
	sess.device.State = model.OnlineState
	sess.device.GatewayIp = config.SystemConfig.LocalIp
	sess.conn = conn

	conn.SetSession(sess)

	its.sessionManager.Add(sess.device.UserId, sess)

	err = its.storageRpc.SaveDevice(sess.device)
	if err != nil {
		_ = sess.conn.Close()
		log.Error("Error saving device", zap.Error(err))
		return
	}
	atomic.AddInt64(&its.clientTotal, 1)

	log.Debug("Device connected successfully", zap.String("userId", sess.device.UserId), zap.String("deviceId", sess.device.DeviceId), zap.String("version", sess.device.DeviceVersion))
}
