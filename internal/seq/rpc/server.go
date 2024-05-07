package rpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/database"
	"eim/pkg/idgenerator"
	"eim/pkg/log"
)

const (
	basePath    = "/eim_seq"
	servicePath = "id"
)

type Config struct {
	Ip                 string
	Port               int
	DatabaseDriver     database.Driver
	DatabaseConnection string
	DatabaseName       string
	EtcdEndpoints      []string
	RedisEndpoints     []string
	RedisPassword      string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

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

	idgenerator.Init(cfg.RedisEndpoints, cfg.RedisPassword)

	db, err := database.NewDatabase(cfg.DatabaseDriver, cfg.DatabaseConnection, cfg.DatabaseName)
	if err != nil {
		log.Error("Error connecting database", zap.String("endpoint", cfg.DatabaseConnection), zap.Error(err))
		return err
	}

	err = svr.RegisterName(servicePath, &segmentSeq{db: db, cache: sync.Map{}}, "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	return err
}
