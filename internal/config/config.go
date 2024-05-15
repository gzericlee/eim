package config

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/util/log"
	"eim/util/net"
)

var SystemConfig *systemConfig

func init() {
	ip, err := net.GetLocalIPv4()
	if err != nil {
		log.Panic("get local ip error -> %v", zap.Error(err))
	}
	SystemConfig = &systemConfig{
		LocalIp: ip,
	}
}

type systemConfig struct {
	LogLevel string
	LocalIp  string
	Redis    struct {
		Endpoints cli.StringSlice
		Password  string
	}
	Etcd struct {
		Endpoints cli.StringSlice
	}
	GatewaySvr struct {
		HttpPort       int
		WebSocketPorts cli.StringSlice
	}
	Mock struct {
		EimEndpoints cli.StringSlice
		ClientCount  int
		MessageCount int
	}
	Database struct {
		Driver     string
		Name       string
		Connection string
	}
	Mq struct {
		Endpoints cli.StringSlice
	}
	SeqSvr struct {
		RpcPort int
	}
	StorageSvr struct {
		RpcPort int
	}
	AuthSvr struct {
		RpcPort int
		Mode    string
	}
	ApiSvr struct {
		HttpPort int
	}
}
