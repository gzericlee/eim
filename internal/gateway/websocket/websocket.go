package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/panjf2000/ants"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eim/internal/model"
	"eim/util/log"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/config"
	"eim/internal/mq"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
)

type Server struct {
	ip             string
	ports          []string
	keepaliveTime  time.Duration
	sessionManager *manager
	workerPool     *ants.Pool
	http           *nbhttp.Server
	seqRpc         *seqrpc.Client
	authRpc        *authrpc.Client
	storageRpc     *storagerpc.Client
	producer       mq.Producer

	receivedMsgTotal int64
	sendMsgTotal     int64
	invalidMsgTotal  int64
	ackTotal         int64
	heartbeatTotal   int64
	clientTotal      int64
	errorTotal       int64
}

func NewServer(ip string, ports []string, seqRpc *seqrpc.Client, authRpc *authrpc.Client, storageRpc *storagerpc.Client, producer mq.Producer) (*Server, error) {
	var address []string
	for _, port := range ports {
		address = append(address, fmt.Sprintf("%v:%v", ip, port))
	}

	taskPool, err := ants.NewPoolPreMalloc(1024)
	if err != nil {
		return nil, fmt.Errorf("new worker pool -> %w", err)
	}

	keepaliveTime := time.Minute * 5

	server := &Server{
		ip:             ip,
		ports:          ports,
		keepaliveTime:  keepaliveTime,
		sessionManager: new(manager),
		workerPool:     taskPool,
		seqRpc:         seqRpc,
		authRpc:        authRpc,
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

	return server, nil
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

	ws := websocket.NewUpgrader()

	ws.OnMessage(its.receiverHandler)
	ws.OnClose(its.closeHandler)
	ws.SetPingHandler(its.pingHandler)

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
		log.Error("Error check auth token", zap.Error(err))
		return
	}

	sess := &session{device: &model.Device{}}
	sess.server = its
	sess.device.UserId = user.BizId
	sess.device.OnlineAt = timestamppb.Now()
	sess.device.DeviceId = r.Header.Get("DeviceId")
	sess.device.DeviceVersion = r.Header.Get("DeviceVersion")
	sess.device.DeviceType = r.Header.Get("DeviceType")
	sess.device.State = model.OnlineState
	sess.device.GatewayIp = config.SystemConfig.LocalIp
	sess.conn = conn
	sess.user = user

	conn.SetSession(sess)

	its.sessionManager.Add(sess.device.UserId, sess)

	err = its.storageRpc.SaveDevice(sess.device)
	if err != nil {
		_ = sess.conn.Close()
		log.Error("Error save device", zap.Error(err))
		return
	}
	atomic.AddInt64(&its.clientTotal, 1)

	sess.sendOfflineMessage()

	log.Debug("device connected successfully", zap.String("userId", sess.device.UserId), zap.String("deviceId", sess.device.DeviceId), zap.String("version", sess.device.DeviceVersion))
}

func (its *Server) closeHandler(conn *websocket.Conn, err error) {
	sess := conn.Session().(*session)
	if sess == nil {
		return
	}

	defer func() {
		atomic.AddInt64(&its.clientTotal, -1)
		its.sessionManager.Remove(sess.device.UserId, sess.device.DeviceId)
		log.Debug("device disconnected", zap.String("deviceId", sess.device.DeviceId), zap.Error(err))
	}()

	sess.device.OfflineAt = timestamppb.Now()
	sess.device.State = model.OfflineState

	err = its.storageRpc.SaveDevice(sess.device)
	if err != nil {
		log.Error("Error save device", zap.Error(err))
		return
	}
}

func (its *Server) pingHandler(conn *websocket.Conn, s string) {
	err := conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))
	if err != nil {
		log.Error("Error set read deadline", zap.Error(err))
		_ = conn.Close()
		return
	}
	err = conn.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
	if err != nil {
		log.Error("Error send pong message", zap.Error(err))
		_ = conn.Close()
		return
	}
	atomic.AddInt64(&its.heartbeatTotal, 1)
}
