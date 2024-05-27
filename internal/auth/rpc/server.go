package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"

	storagerpc "eim/internal/storage/rpc"
)

const (
	basePath    = "/eim_auth"
	servicePath = "auth"
)

type Config struct {
	Ip            string
	Port          int
	EtcdEndpoints []string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

	storageRpc, err := storagerpc.NewClient(cfg.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("new storage rpc client -> %w", err)
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

	err = svr.RegisterName(servicePath, &Authentication{StorageRpc: storageRpc}, "")
	if err != nil {
		return fmt.Errorf("register auth service -> %w", err)
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	if err != nil {
		return fmt.Errorf("start server -> %w", err)
	}

	return nil
}
