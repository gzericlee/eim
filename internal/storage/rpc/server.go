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
	"eim/pkg/lock"
)

const (
	basePath = "eim_storage"

	deviceServicePath    = "device"
	messageServicePath   = "message"
	bizServicePath       = "biz"
	bizMemberServicePath = "biz_member"
	gatewayServicePath   = "gateway"
	segmentServicePath   = "segment"
	refresherServicePath = "refresher"

	bizCachePool       = "bizs"
	deviceCachePool    = "devices"
	bizMemberCachePool = "biz_members"

	cacheKeyFormat = "%s:%s:%s"
)

var (
	storageRpc  *Client
	singleGroup singleflight.Group
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

	svr.AsyncWrite = true

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

	deviceCache, err := cache.NewCache("device", 3*1024*1024*1024, 1000000, true)
	if err != nil {
		panic(err)
	}

	bizCache, err := cache.NewCache("biz", 3*1024*1024*1024, 500000, true)
	if err != nil {
		panic(err)
	}

	bizMemberCache, err := cache.NewCache("biz_member", 3*1024*1024*1024, 1000000, true)
	if err != nil {
		panic(err)
	}

	keyLock := lock.NewKeyLock()

	err = svr.RegisterName(refresherServicePath, &Refresher{deviceCache: deviceCache, bizCache: bizCache, bizMemberCache: bizMemberCache, lock: keyLock}, "")
	if err != nil {
		return fmt.Errorf("register service(%s) -> %w", refresherServicePath, err)
	}

	for _, service := range cfg.RegistryServices {
		var rcvr interface{}
		switch service {
		case deviceServicePath:
			{
				rcvr = &Device{storageCache: deviceCache, database: db, redisManager: redisManager, lock: keyLock}
			}
		case messageServicePath:
			{
				rcvr = &Message{database: db, redisManager: redisManager}
			}
		case bizServicePath:
			{
				rcvr = &Biz{storageCache: bizCache, database: db, redisManager: redisManager}
			}
		case bizMemberServicePath:
			{
				rcvr = &BizMember{storageCache: bizMemberCache, redisManager: redisManager}
			}
		case gatewayServicePath:
			{
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

	storageRpc, err = NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("new storage rpc client -> %w", err)
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
