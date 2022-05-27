package storage

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/global"
)

func InitStorageServer(ip string, port int, etcdEndpoints []string) error {
	svr := server.NewServer()

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", ip, port),
		EtcdServers:    etcdEndpoints,
		BasePath:       "/eim_storage",
		UpdateInterval: time.Minute,
	}
	err := plugin.Start()
	if err != nil {
		global.Logger.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName("Device", new(Device), "")
	if err != nil {
		return err
	}

	err = svr.RegisterName("Message", new(Message), "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
