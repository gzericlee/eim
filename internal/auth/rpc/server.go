package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/redis"
	"eim/pkg/log"
)

const (
	basePath    = "/eim_auth"
	servicePath = "auth"
)

type Config struct {
	Ip             string
	Port           int
	EtcdEndpoints  []string
	RedisEndpoints []string
	RedisPassword  string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

	redisManager, err := redis.NewManager(cfg.RedisEndpoints, cfg.RedisPassword)
	if err != nil {
		log.Error("Error connection redis cluster", zap.Error(err))
		return err
	}

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err = plugin.Start()
	if err != nil {
		log.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName(servicePath, &Authentication{RedisManager: redisManager}, "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	return err
}
