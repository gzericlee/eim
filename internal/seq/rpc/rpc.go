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

	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	"github.com/gzericlee/eim/internal/seq/rpc/service"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	eimmetrics "github.com/gzericlee/eim/pkg/metrics"
	"github.com/gzericlee/eim/pkg/snowflake"
)

const (
	basePath    = "github.com/gzericlee/eim_seq"
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

	segmentRpc, err := storagerpc.NewSegmentClient(cfg.EtcdEndpoints)
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

	err = svr.RegisterName(servicePath, service.NewSeqService(segmentRpc, generator), "")
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

func NewSeqClient(etcdEndpoints []string) (*seqrpc.SeqClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, servicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for seq service -> %w", err)
	}
	return &seqrpc.SeqClient{
		XClientPool: rpcxclient.NewXClientPool(100, servicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}
