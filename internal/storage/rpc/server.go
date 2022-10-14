package rpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/database/maindb"
	"eim/pkg/log"
)

const (
	basePath     = "/eim_storage"
	servicePath1 = "Device"
	servicePath2 = "Message"
)

var mainDb = maindb.NewMainDB()

func StartServer(ip string, port int, etcdEndpoints []string) error {
	svr := server.NewServer()

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", ip, port),
		EtcdServers:    etcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err := plugin.Start()
	if err != nil {
		log.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName(servicePath1, new(Device), "")
	if err != nil {
		return err
	}

	err = svr.RegisterName(servicePath2, new(Message), "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
