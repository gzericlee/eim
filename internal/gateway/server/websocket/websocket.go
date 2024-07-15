package websocket

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/panjf2000/ants"

	authrpc "github.com/gzericlee/eim/internal/auth/rpc/client"
	"github.com/gzericlee/eim/internal/gateway/session"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/mq"
	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type Server struct {
	ip             string
	port           int
	keepaliveTime  time.Duration
	sessionManager *session.Manager
	workerPool     *ants.Pool
	http           *nbhttp.Server
	seqRpc         *seqrpc.SeqClient
	authRpc        *authrpc.AuthClient
	deviceRpc      *storagerpc.DeviceClient
	messageRpc     *storagerpc.MessageClient
	gatewayRpc     *storagerpc.GatewayClient
	producer       mq.IProducer

	receivedMsgTotal int64
	sendMsgTotal     int64
	invalidMsgTotal  int64
	ackTotal         int64
	heartbeatTotal   int64
	clientTotal      int64
	errorTotal       int64
}

type Config struct {
	Ip         string
	Port       int
	SeqRpc     *seqrpc.SeqClient
	AuthRpc    *authrpc.AuthClient
	DeviceRpc  *storagerpc.DeviceClient
	MessageRpc *storagerpc.MessageClient
	GatewayRpc *storagerpc.GatewayClient
	Producer   mq.IProducer
}

func NewServer(cfg *Config) (*Server, error) {
	taskPool, err := ants.NewPoolPreMalloc(1024)
	if err != nil {
		return nil, fmt.Errorf("new worker pool -> %w", err)
	}

	keepaliveTime := time.Minute * 5

	server := &Server{
		ip:             cfg.Ip,
		port:           cfg.Port,
		keepaliveTime:  keepaliveTime,
		sessionManager: session.NewManager(),
		workerPool:     taskPool,
		seqRpc:         cfg.SeqRpc,
		authRpc:        cfg.AuthRpc,
		deviceRpc:      cfg.DeviceRpc,
		messageRpc:     cfg.MessageRpc,
		gatewayRpc:     cfg.GatewayRpc,
		producer:       cfg.Producer,
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", server.connect)

	server.http = nbhttp.NewServer(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port)},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
		KeepaliveTime:           keepaliveTime,
	})

	return server, nil
}

func (its *Server) Start() error {
	return its.http.Start()
}

func (its *Server) Stop() {
	its.http.Stop()
}

func (its *Server) GetGatewayRpc() *storagerpc.GatewayClient {
	return its.gatewayRpc
}

func (its *Server) GetMessageRpc() *storagerpc.MessageClient {
	return its.messageRpc
}

func (its *Server) GetDeviceRpc() *storagerpc.DeviceClient {
	return its.deviceRpc
}

func (its *Server) GetSeqRpc() *seqrpc.SeqClient {
	return its.seqRpc
}

func (its *Server) GetAuthRpc() *authrpc.AuthClient {
	return its.authRpc
}

func (its *Server) GetMQProducer() mq.IProducer {
	return its.producer
}

func (its *Server) GetSessionManager() *session.Manager {
	return its.sessionManager
}

func (its *Server) IncrReceivedMsgTotal(count int64) {
	atomic.AddInt64(&its.receivedMsgTotal, count)
}

func (its *Server) IncrSendMsgTotal(count int64) {
	atomic.AddInt64(&its.sendMsgTotal, count)
}

func (its *Server) IncrInvalidMsgTotal(count int64) {
	atomic.AddInt64(&its.invalidMsgTotal, count)
}

func (its *Server) IncrAckTotal(count int64) {
	atomic.AddInt64(&its.ackTotal, count)
}

func (its *Server) IncrHeartbeatTotal(count int64) {
	atomic.AddInt64(&its.heartbeatTotal, count)
}

func (its *Server) IncrClientTotal(count int64) {
	atomic.AddInt64(&its.clientTotal, count)
}

func (its *Server) IncrErrorTotal(count int64) {
	atomic.AddInt64(&its.errorTotal, count)
}

func (its *Server) RegistryGateway() {
	err := its.gatewayRpc.RegisterGateway(&model.Gateway{
		Ip:            its.ip,
		Port:          int32(its.port),
		ClientTotal:   its.clientTotal,
		SendTotal:     its.sendMsgTotal,
		ReceivedTotal: its.receivedMsgTotal,
		InvalidTotal:  its.invalidMsgTotal,
	}, time.Second*10)
	if err != nil {
		return
	}
}
