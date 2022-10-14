package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/types"
	"eim/pkg/log"
)

const (
	basePath    = "/eim_auth"
	servicePath = "Auth"
)

type Request struct {
	Token string
}

type Reply struct {
	User *types.User
}

type Authentication struct {
}

func (its *Authentication) CheckToken(ctx context.Context, req *Request, reply *Reply) error {
	authenticator := newAuthenticator(mode(config.SystemConfig.AuthSvr.Mode))
	user, err := authenticator.CheckToken(req.Token)
	reply.User = user
	return err
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

	err = svr.RegisterName(servicePath, new(Authentication), "")
	if err != nil {
		return err
	}

	err = svr.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
