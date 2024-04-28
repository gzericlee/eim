package gateway

import (
	"github.com/lesismal/nbio/logging"
	"go.uber.org/zap"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/gateway/websocket"
	"eim/internal/mq"
	"eim/internal/redis"
	seqrpc "eim/internal/seq/rpc"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
)

type Config struct {
	Ip             string
	Ports          []string
	NsqEndpoints   []string
	EtcdEndpoints  []string
	RedisEndpoints []string
	RedisPassword  string
}

func StartWebsocketServer(cfg Config) (*websocket.Server, error) {
	logging.SetLevel(logging.LevelNone)

	seqRpc, err := seqrpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, err
	}

	authRpc, err := authrpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, err
	}

	storageRpc, err := storagerpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return nil, err
	}

	redisManager, err := redis.NewManager(cfg.RedisEndpoints, cfg.RedisPassword)
	if err != nil {
		return nil, err
	}

	producer, err := mq.NewProducer(cfg.NsqEndpoints)
	if err != nil {
		return nil, err
	}

	server, err := websocket.NewServer(cfg.Ip, cfg.Ports, seqRpc, authRpc, storageRpc, redisManager, producer)
	if err != nil {
		return nil, err
	}

	err = server.Start()
	if err != nil {
		return nil, err
	}

	log.Info("Listening websocket connect", zap.String("ip", cfg.Ip), zap.Strings("ports", cfg.Ports))

	return server, nil
}
