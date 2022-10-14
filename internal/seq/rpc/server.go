package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/pkg/log"
)

const (
	basePath    = "/eim_seq"
	servicePath = "Id"
)

type Request struct {
	Id string
}

type Reply struct {
	Number int64
}

type Seq struct {
}

func (its *Seq) Number(ctx context.Context, req *Request, reply *Reply) error {
	reply.Number = newSeq(req.Id).Get()
	return nil
}

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

	err = svr.RegisterName(servicePath, new(Seq), "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
