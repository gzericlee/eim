package gateway

import (
	"fmt"

	"github.com/lesismal/nbio/logging"
	"go.uber.org/zap"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/gateway/websocket"
	"eim/internal/mq"
	"eim/internal/redis"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

type Config struct {
	Ip             string
	Ports          []string
	MqEndpoints    []string
	EtcdEndpoints  []string
	RedisEndpoints []string
	RedisPassword  string
}

func StartWebsocketServer(cfg Config) (*websocket.Server, error) {
	logging.SetLevel(logging.LevelNone)

	seqRpc, err := seqrpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new seq rpc client -> %w", err)
	}

	authRpc, err := authrpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new auth rpc client -> %w", err)
	}

	storageRpc, err := storagerpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new storage rpc client -> %w", err)
	}

	redisManager, err := redis.NewManager(redis.Config{
		RedisEndpoints: cfg.RedisEndpoints,
		RedisPassword:  cfg.RedisPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("new redis manager -> %w", err)
	}

	producer, err := mq.NewProducer(cfg.MqEndpoints)
	if err != nil {
		return nil, fmt.Errorf("new mq producer -> %w", err)
	}

	server, err := websocket.NewServer(cfg.Ip, cfg.Ports, seqRpc, authRpc, storageRpc, redisManager, producer)
	if err != nil {
		return nil, fmt.Errorf("new websocket server -> %w", err)
	}

	err = server.Start()
	if err != nil {
		return nil, fmt.Errorf("start websocket server -> %w", err)
	}

	log.Info("Listening websocket connect", zap.String("ip", cfg.Ip), zap.Strings("ports", cfg.Ports))

	return server, nil
}
