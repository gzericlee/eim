package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/database"
	"eim/internal/redis"
	"eim/pkg/log"
)

const (
	basePath     = "/eim_storage"
	servicePath1 = "device"
	servicePath2 = "message"
)

type Config struct {
	Ip                 string
	Port               int
	EtcdEndpoints      []string
	DatabaseDriver     database.Driver
	DatabaseConnection string
	DatabaseName       string
	RedisEndpoints     []string
	RedisPassword      string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

	db, err := database.NewDatabase(cfg.DatabaseDriver, cfg.DatabaseConnection, cfg.DatabaseName)
	if err != nil {
		log.Error("Error connecting database", zap.String("endpoint", cfg.DatabaseConnection), zap.Error(err))
		return err
	}

	redisManager, err := redis.NewManager(cfg.RedisEndpoints, cfg.RedisPassword)
	if err != nil {
		log.Error("Error connecting redis cluster", zap.Strings("endpoints", cfg.RedisEndpoints), zap.Error(err))
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

	err = svr.RegisterName(servicePath1, &Device{RedisManager: redisManager, Database: db}, "")
	if err != nil {
		return err
	}

	err = svr.RegisterName(servicePath2, &Message{Database: db}, "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	return err
}
