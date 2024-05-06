package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/redis"
	"eim/internal/seq/medis"
	"eim/pkg/log"
)

const (
	basePath    = "/eim_seq"
	servicePath = "id"
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

	//client, err := clientv3.New(clientv3.Config{
	//	Endpoints:   cfg.EtcdEndpoints,
	//	DialTimeout: 5 * time.Second,
	//})
	//if err != nil {
	//	log.Error("Error connecting etcd cluster", zap.Error(err))
	//	return err
	//}

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err := plugin.Start()
	if err != nil {
		log.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	redisManager, err := redis.NewManager(cfg.RedisEndpoints, cfg.RedisPassword)
	if err != nil {
		log.Error("Error connecting redis cluster", zap.Strings("endpoints", cfg.RedisEndpoints), zap.Error(err))
		return err
	}

	instance, err := medis.NewInstance(redisManager.GetRedisClient())
	if err != nil {
		log.Error("Error creating medis instance", zap.Error(err))
		return err
	}

	err = svr.RegisterName(servicePath, &medisSeq{instance: instance}, "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	return err
}
