package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"golang.org/x/sync/singleflight"

	"eim/internal/database"
	"eim/internal/redis"
	"eim/pkg/cache"
	"eim/pkg/cache/notify"
	"eim/pkg/cache/stats"
	"eim/util/log"
)

const (
	basePath             = "/eim_storage"
	deviceServicePath    = "device"
	messageServicePath   = "message"
	userServicePath      = "user"
	bizMemberServicePath = "biz_member"
	gatewayServicePath   = "gateway"
	segmentServicePath   = "segment"

	userCachePool      = "users"
	gatewayCachePool   = "gateways"
	deviceCachePool    = "devices"
	bizMemberCachePool = "biz_members"
)

var (
	group singleflight.Group
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
	RegistryServices     []string
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

	storageCache := cache.NewLRUCache(1024, 65535, 0)

	notify.Init(notify.NewRedisNotifier(redisManager.GetRedisClient()))

	for _, service := range cfg.RegistryServices {
		var rcvr interface{}
		switch service {
		case deviceServicePath:
			{
				err = notify.Bind(deviceCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("notify bind(device) -> %w", err)
				}
				err = stats.Bind(deviceCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("stats bind(device) -> %w", err)
				}
				rcvr = &Device{storageCache: storageCache, database: db}
			}
		case messageServicePath:
			{
				rcvr = &Message{database: db, redisManager: redisManager}
			}
		case userServicePath:
			{
				err = notify.Bind(userCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("notify bind(user) -> %w", err)
				}
				err = stats.Bind(userCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("stats bind(user) -> %w", err)
				}
				rcvr = &User{storageCache: storageCache, database: db}
			}
		case bizMemberServicePath:
			{
				err = notify.Bind(bizMemberCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("notify bind(biz_member) -> %w", err)
				}
				err = stats.Bind(bizMemberCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("stats bind(biz_member) -> %w", err)
				}
				rcvr = &BizMember{storageCache: storageCache, redisManager: redisManager}
			}
		case gatewayServicePath:
			{
				err = notify.Bind(gatewayCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("notify bind(gateway) -> %w", err)
				}
				err = stats.Bind(gatewayCachePool, storageCache)
				if err != nil {
					return fmt.Errorf("stats bind(gateway) -> %w", err)
				}
				rcvr = &Gateway{redisManager: redisManager}
			}
		case segmentServicePath:
			{
				rcvr = &Segment{database: db}
			}
		}
		err = svr.RegisterName(service, rcvr, "")
		if err != nil {
			return fmt.Errorf("register service(%s) -> %w", service, err)
		}
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			stats.All().Range(func(k, v interface{}) bool {
				log.Info(fmt.Sprintf("Cache stats: %s %+v", k, v))
				return true
			})
		}
	}()

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
