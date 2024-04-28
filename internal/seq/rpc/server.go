package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"eim/pkg/log"
)

const (
	basePath    = "/eim_seq"
	servicePath = "id"
)

type Config struct {
	Ip            string
	Port          int
	EtcdEndpoints []string
}

func StartServer(cfg Config) error {
	svr := server.NewServer()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("Error connecting etcd cluster", zap.Error(err))
		return err
	}

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", cfg.Ip, cfg.Port),
		EtcdServers:    cfg.EtcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err = plugin.Start()
	if err != nil {
		log.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName(servicePath, &Seq{etcdClient: client}, "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port))
	return err
}
