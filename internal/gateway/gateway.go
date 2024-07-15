package gateway

import (
	"fmt"
	"time"

	"github.com/lesismal/nbio/logging"
	"go.uber.org/zap"

	authrpc "github.com/gzericlee/eim/internal/auth/rpc"
	"github.com/gzericlee/eim/internal/gateway/server"
	"github.com/gzericlee/eim/internal/gateway/server/websocket"
	"github.com/gzericlee/eim/internal/mq"
	seqrpc "github.com/gzericlee/eim/internal/seq/rpc"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	"github.com/gzericlee/eim/pkg/log"
)

type Config struct {
	Ip            string
	Port          int
	MqEndpoints   []string
	EtcdEndpoints []string
}

func StartWebsocketServer(cfg *Config) (server.IServer, error) {
	logging.SetLevel(logging.LevelNone)

	seqRpc, err := seqrpc.NewSeqClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new seq rpc client -> %w", err)
	}

	authRpc, err := authrpc.NewAuthClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new auth rpc client -> %w", err)
	}

	deviceRpc, err := storagerpc.NewDeviceClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new device rpc client -> %w", err)
	}

	messageRpc, err := storagerpc.NewMessageClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new message rpc client -> %w", err)
	}

	gatewayRpc, err := storagerpc.NewGatewayClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new gateway rpc client -> %w", err)
	}

	producer, err := mq.NewProducer(cfg.MqEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new mq producer -> %w", err)
	}

	server, err := websocket.NewServer(&websocket.Config{
		Ip:         cfg.Ip,
		Port:       cfg.Port,
		SeqRpc:     seqRpc,
		AuthRpc:    authRpc,
		DeviceRpc:  deviceRpc,
		MessageRpc: messageRpc,
		GatewayRpc: gatewayRpc,
		Producer:   producer,
	})
	if err != nil {
		return nil, fmt.Errorf("new websocket server -> %w", err)
	}

	err = server.Start()
	if err != nil {
		return nil, fmt.Errorf("start websocket server -> %w", err)
	}

	go server.PrintServiceStats()

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				{
					server.RegistryGateway()
				}
			}
		}
	}()

	log.Info("Listening websocket connect", zap.String("ip", cfg.Ip), zap.Int("port", cfg.Port))

	return server, nil
}
