package server

import (
	authrpc "github.com/gzericlee/eim/internal/auth/rpc/client"
	"github.com/gzericlee/eim/internal/gateway/session"
	"github.com/gzericlee/eim/internal/mq"
	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type IServer interface {
	Start() error
	Stop()

	Send(sess *session.Session, cmd int, body []byte)

	GetMessageRpc() *storagerpc.MessageClient
	GetGatewayRpc() *storagerpc.GatewayClient
	GetDeviceRpc() *storagerpc.DeviceClient
	GetSeqRpc() *seqrpc.SeqClient
	GetAuthRpc() *authrpc.AuthClient
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
