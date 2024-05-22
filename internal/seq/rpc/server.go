package rpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"

	"eim/internal/database"
	"eim/pkg/snowflake"
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
		return fmt.Errorf("start etcd v3 register plugin -> %w", err)
	}
	svr.Plugins.Add(plugin)

	generator, err := snowflake.NewGenerator(snowflake.GeneratorConfig{
		RedisEndpoints: cfg.RedisEndpoints,
		RedisPassword:  cfg.RedisPassword,
		MaxWorkerId:    1023,
		MinWorkerId:    1,
		NodeCount:      5,
	})
	if err != nil {
		return fmt.Errorf("new snowflake id incrementer -> %w", err)
	}

	db, err := database.NewDatabase(cfg.DatabaseDriver, cfg.DatabaseConnection, cfg.DatabaseName)
	if err != nil {
		return fmt.Errorf("new database -> %w", err)
	}

	err = svr.RegisterName(servicePath, &segmentSeq{db: db, cache: sync.Map{}, generator: generator}, "")
	if err != nil {
		return fmt.Errorf("register seq service -> %w", err)
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
