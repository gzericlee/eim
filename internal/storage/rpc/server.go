package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"

	"eim/internal/database"
	"eim/internal/redis"
)

const (
	basePath     = "/eim_storage"
	servicePath1 = "device"
	servicePath2 = "message"
)

type Config struct {
	Ip                   string
	Port                 int
	EtcdEndpoints        []string
	DatabaseDriver       database.Driver
	DatabaseConnection   string
	DatabaseName         string
	RedisEndpoints       []string
	RedisPassword        string
	OfflineMessageExpire int
	OfflineDeviceExpire  int
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

	db, err := database.NewDatabase(cfg.DatabaseDriver, cfg.DatabaseConnection, cfg.DatabaseName)
	if err != nil {
		return fmt.Errorf("new database -> %w", err)
	}

	redisManager, err := redis.NewManager(redis.Config{
		RedisEndpoints:       cfg.RedisEndpoints,
		RedisPassword:        cfg.RedisPassword,
		OfflineMessageExpire: time.Hour * 24 * time.Duration(cfg.OfflineMessageExpire),
		OfflineDeviceExpire:  time.Hour * 24 * time.Duration(cfg.OfflineDeviceExpire),
	})
	if err != nil {
		return fmt.Errorf("new redis manager -> %w", err)
	}

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err = plugin.Start()
	if err != nil {
		return fmt.Errorf("start etcd v3 register plugin -> %w", err)
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName(servicePath1, &Device{RedisManager: redisManager, Database: db}, "")
	if err != nil {
		return fmt.Errorf("register device service -> %w", err)
	}

	err = svr.RegisterName(servicePath2, &Message{Database: db}, "")
	if err != nil {
		return fmt.Errorf("register message service -> %w", err)
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
