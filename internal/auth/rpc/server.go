package rpc

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	etcdplugin "github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	rpcxplugin "github.com/smallnest/rpcx/serverplugin"

	storagerpc "eim/internal/storage/rpc"
	eimmetrics "eim/pkg/metrics"
)

const (
	basePath    = "eim_auth"
	servicePath = "auth"
)

type Config struct {
	Ip            string
	Port          int
	EtcdEndpoints []string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()
	svr.AsyncWrite = true

	storageRpc, err := storagerpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("new storage rpc client -> %w", err)
	}

	metricsPlugin := rpcxplugin.NewMetricsPlugin(metrics.DefaultRegistry)

	etcdPlugin := &etcdplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err = etcdPlugin.Start()
	if err != nil {
		return fmt.Errorf("start etcd v3 register plugin -> %w", err)
	}

	svr.Plugins.Add(etcdPlugin)
	svr.Plugins.Add(metricsPlugin)

	err = svr.RegisterName(servicePath, &Authentication{StorageRpc: storageRpc}, "")
	if err != nil {
		return fmt.Errorf("register auth service -> %w", err)
	}

	eimmetrics.EnableMetrics(32003)

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
