package rpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	etcdplugin "github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	rpcxplugin "github.com/smallnest/rpcx/serverplugin"

	storagerpc "eim/internal/storage/rpc"
	eimmetrics "eim/pkg/metrics"
	"eim/pkg/snowflake"
)

const (
	basePath    = "eim_seq"
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
	svr.AsyncWrite = true

	metricsPlugin := rpcxplugin.NewMetricsPlugin(metrics.DefaultRegistry)

	etcdPlugin := &etcdplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err := etcdPlugin.Start()
	if err != nil {
		return fmt.Errorf("start etcd v3 register plugin -> %w", err)
	}

	svr.Plugins.Add(etcdPlugin)
	svr.Plugins.Add(metricsPlugin)

	storageRpc, err := storagerpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("new storage rpc client -> %w", err)
	}

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

	err = svr.RegisterName(servicePath, &segmentSeq{storageRpc: storageRpc, cache: sync.Map{}, generator: generator}, "")
	if err != nil {
		return fmt.Errorf("register seq service -> %w", err)
	}

	eimmetrics.EnableMetrics(32003)

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
