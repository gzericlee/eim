package seq

import (
	"context"
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/global"
)

const (
	basePath    = "/eim_seq"
	servicePath = "Id"
)

type Request struct {
	UserId string
}

type Reply struct {
	Id int64
}

type Seq int

func (its *Seq) Id(ctx context.Context, req *Request, reply *Reply) error {
	reply.Id = newId(req.UserId).Get()
	return nil
}

func InitSeqServer(ip string, port int, etcdEndpoints []string) error {
	svr := server.NewServer()

	plugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", ip, port),
		EtcdServers:    etcdEndpoints,
		BasePath:       basePath,
		UpdateInterval: time.Minute,
	}
	err := plugin.Start()
	if err != nil {
		global.Logger.Error("Error registering etcd plugin", zap.Error(err))
		return err
	}
	svr.Plugins.Add(plugin)

	err = svr.RegisterName(servicePath, new(Seq), "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
