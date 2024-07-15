package rpc

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	rpcxetcdclient "github.com/rpcxio/rpcx-etcd/client"
	etcdplugin "github.com/rpcxio/rpcx-etcd/serverplugin"
	rpcxclient "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	rpcxplugin "github.com/smallnest/rpcx/serverplugin"

	"github.com/gzericlee/eim/internal/auth"
	"github.com/gzericlee/eim/internal/auth/rpc/client"
	"github.com/gzericlee/eim/internal/auth/rpc/service"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	eimmetrics "github.com/gzericlee/eim/pkg/metrics"
	oauth2lib "github.com/gzericlee/eim/pkg/oauth2"
)

const (
	basePath        = "github.com/gzericlee/eim_auth"
	authServicePath = "auth"
)

type Config struct {
	Ip            string
	Port          int
	EtcdEndpoints []string
	AuthMode      auth.Mode
	Oauth2Client  oauth2lib.Client
}

func StartServer(cfg Config) error {
	svr := server.NewServer()
	svr.AsyncWrite = true

	bizRpc, err := storagerpc.NewBizClient(cfg.EtcdEndpoints)
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

	err = svr.RegisterName(authServicePath, service.NewAuthService(cfg.AuthMode, cfg.Oauth2Client, bizRpc), "")
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

func NewAuthClient(etcdEndpoints []string) (*client.AuthClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, authServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for auth service -> %w", err)
	}

	return &client.AuthClient{
		XClientPool: rpcxclient.NewXClientPool(100, authServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}
