package websocket

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/panjf2000/ants"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/gateway/session"
	"eim/internal/mq"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
)

type Server struct {
	ip             string
	ports          []string
	keepaliveTime  time.Duration
	sessionManager *session.Manager
	workerPool     *ants.Pool
	http           *nbhttp.Server
	seqRpc         *seqrpc.Client
	authRpc        *authrpc.Client
	storageRpc     *storagerpc.Client
	producer       mq.IProducer

	receivedMsgTotal int64
	sendMsgTotal     int64
	invalidMsgTotal  int64
	ackTotal         int64
	heartbeatTotal   int64
	clientTotal      int64
	errorTotal       int64
}

func NewServer(ip string, ports []string, seqRpc *seqrpc.Client, authRpc *authrpc.Client, storageRpc *storagerpc.Client, producer mq.IProducer) (*Server, error) {
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
		sessionManager: session.NewManager(),
		workerPool:     taskPool,
		seqRpc:         seqRpc,
		authRpc:        authRpc,
		storageRpc:     storageRpc,
		producer:       producer,
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", server.connect)

	server.http = nbhttp.NewServer(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   address,
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

func (its *Server) GetStorageRpc() *storagerpc.Client {
	return its.storageRpc
}

func (its *Server) GetSeqRpc() *seqrpc.Client {
	return its.seqRpc
}

func (its *Server) GetAuthRpc() *authrpc.Client {
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
