package config

import (
	"github.com/urfave/cli/v2"

	"eim/util"
)

var SystemConfig *systemConfig

func init() {
	ip, err := util.GetLocalIpV4()
	if err != nil {
		panic(err)
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
	Nsq struct {
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
