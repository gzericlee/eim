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

	"github.com/gzericlee/eim/internal/database"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/redis"
	"github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/internal/storage/rpc/service"
	"github.com/gzericlee/eim/pkg/cache"
	"github.com/gzericlee/eim/pkg/lock"
	eimmetrics "github.com/gzericlee/eim/pkg/metrics"
)

const (
	basePath = "github.com/gzericlee/eim_storage"

	deviceServicePath    = "device"
	messageServicePath   = "message"
	bizServicePath       = "biz"
	bizMemberServicePath = "biz_member"
	gatewayServicePath   = "gateway"
	segmentServicePath   = "segment"
	refresherServicePath = "refresher"
	tenantServicePath    = "tenant"
	fileServicePath      = "file"
	fileThumbServicePath = "file_thumb"
)

var (
	refreshRpc *client.RefresherClient
)

type Config struct {
	Ip                   string
	Port                 int
	EtcdEndpoints        []string
	DatabaseDriver       database.Driver
	DatabaseConnections  []string
	DatabaseName         string
	RedisEndpoints       []string
	RedisPassword        string
	OfflineMessageExpire int
	OfflineDeviceExpire  int
	RegistryServices     []string
}

func StartServer(cfg *Config) error {
	svr := server.NewServer()
	svr.AsyncWrite = true

	db, err := database.NewDatabase(cfg.DatabaseDriver, cfg.DatabaseConnections, cfg.DatabaseName)
	if err != nil {
		return fmt.Errorf("new database -> %w", err)
	}

	redisManager, err := redis.NewManager(&redis.Config{
		RedisEndpoints:       cfg.RedisEndpoints,
		RedisPassword:        cfg.RedisPassword,
		OfflineMessageExpire: time.Hour * 24 * time.Duration(cfg.OfflineMessageExpire),
		OfflineDeviceExpire:  time.Hour * 24 * time.Duration(cfg.OfflineDeviceExpire),
	})
	if err != nil {
		return fmt.Errorf("new redis manager -> %w", err)
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

	devicesCache, err := cache.NewCache[string, []*model.Device]("devices", 1000000)
	if err != nil {
		panic(err)
	}

	bizCache, err := cache.NewCache[string, *model.Biz]("biz", 1000000)
	if err != nil {
		panic(err)
	}

	bizMembersCache, err := cache.NewCache[string, []string]("biz_members", 1000000)
	if err != nil {
		panic(err)
	}

	tenantCache, err := cache.NewCache[string, *model.Tenant]("tenant", 100000)
	if err != nil {
		panic(err)
	}

	keyLock := lock.NewKeyLock()

	err = svr.RegisterName(refresherServicePath, service.NewRefresherService(keyLock, devicesCache, bizCache, bizMembersCache, tenantCache), "")
	if err != nil {
		return fmt.Errorf("register service(%s) -> %w", refresherServicePath, err)
	}

	refreshRpc, err = NewRefresherClient(cfg.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("new storage rpc client -> %w", err)
	}

	for _, svcName := range cfg.RegistryServices {
		var rcvr interface{}
		switch svcName {
		case deviceServicePath:
			{
				rcvr = service.NewDeviceService(keyLock, devicesCache, db, refreshRpc)
			}
		case messageServicePath:
			{
				rcvr = service.NewMessageService(db, redisManager)
			}
		case bizServicePath:
			{
				rcvr = service.NewBizService(bizCache, db, refreshRpc)
			}
		case bizMemberServicePath:
			{
				rcvr = service.NewBizMemberService(bizMembersCache, db, refreshRpc)
			}
		case gatewayServicePath:
			{
				rcvr = service.NewGatewayService(redisManager)
			}
		case segmentServicePath:
			{
				rcvr = service.NewSegmentService(db)
			}
		case tenantServicePath:
			{
				rcvr = service.NewTenantService(tenantCache, db, refreshRpc)
			}
		case fileServicePath:
			{
				rcvr = service.NewFileService(db)
			}
		case fileThumbServicePath:
			{
				rcvr = service.NewFileThumbService(db)
			}
		}
		err = svr.RegisterName(svcName, rcvr, "")
		if err != nil {
			return fmt.Errorf("register service(%s) -> %w", svcName, err)
		}
	}

	eimmetrics.EnableMetrics(32003)

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}

func NewBizClient(etcdEndpoints []string) (*client.BizClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, bizServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for biz service -> %w", err)
	}
	return &client.BizClient{
		XClientPool: rpcxclient.NewXClientPool(100, bizServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewBizMemberClient(etcdEndpoints []string) (*client.BizMemberClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, bizMemberServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for biz_member service -> %w", err)
	}
	return &client.BizMemberClient{
		XClientPool: rpcxclient.NewXClientPool(100, bizMemberServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewDeviceClient(etcdEndpoints []string) (*client.DeviceClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, deviceServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for device service -> %w", err)
	}
	return &client.DeviceClient{
		XClientPool: rpcxclient.NewXClientPool(100, deviceServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewGatewayClient(etcdEndpoints []string) (*client.GatewayClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, gatewayServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for gateway service -> %w", err)
	}
	return &client.GatewayClient{
		XClientPool: rpcxclient.NewXClientPool(100, gatewayServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewMessageClient(etcdEndpoints []string) (*client.MessageClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, messageServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for message service -> %w", err)
	}
	return &client.MessageClient{
		XClientPool: rpcxclient.NewXClientPool(100, messageServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewTenantClient(etcdEndpoints []string) (*client.TenantClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, tenantServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for tenant service -> %w", err)
	}
	return &client.TenantClient{
		XClientPool: rpcxclient.NewXClientPool(100, tenantServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewFileClient(etcdEndpoints []string) (*client.FileClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, fileServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for file service -> %w", err)
	}
	return &client.FileClient{
		XClientPool: rpcxclient.NewXClientPool(100, fileServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewFileThumbClient(etcdEndpoints []string) (*client.FileThumbClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, fileThumbServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for file_thumb service -> %w", err)
	}
	return &client.FileThumbClient{
		XClientPool: rpcxclient.NewXClientPool(100, fileThumbServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewSegmentClient(etcdEndpoints []string) (*client.SegmentClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, segmentServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for segment service -> %w", err)
	}
	return &client.SegmentClient{
		XClientPool: rpcxclient.NewXClientPool(100, segmentServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}

func NewRefresherClient(etcdEndpoints []string) (*client.RefresherClient, error) {
	discovery, err := rpcxetcdclient.NewEtcdV3Discovery(basePath, refresherServicePath, etcdEndpoints, false, nil)
	if err != nil {
		return nil, fmt.Errorf("new etcd v3 discovery for refresher service -> %w", err)
	}
	return &client.RefresherClient{
		XClientPool: rpcxclient.NewXClientPool(100, refresherServicePath, rpcxclient.Failover, rpcxclient.RoundRobin, discovery, rpcxclient.DefaultOption),
	}, nil
}
