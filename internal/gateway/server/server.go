package server

import (
	authrpc "eim/internal/auth/rpc"
	"eim/internal/gateway/session"
	"eim/internal/mq"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
)

type IServer interface {
	Start() error
	Stop()

	Send(sess *session.Session, cmd int, body []byte)

	GetStorageRpc() *storagerpc.Client
	GetSeqRpc() *seqrpc.Client
	GetAuthRpc() *authrpc.Client
	GetMQProducer() mq.IProducer
	GetSessionManager() *session.Manager

	IncrReceivedMsgTotal(count int64)
	IncrSendMsgTotal(count int64)
	IncrInvalidMsgTotal(count int64)
	IncrAckTotal(count int64)
	IncrHeartbeatTotal(count int64)
	IncrClientTotal(count int64)
	IncrErrorTotal(count int64)

	PrintServiceStats()
	RegistryGateway()
}
