package global

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
	SystemConfig = &systemConfig{LocalIp: ip}
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
	Gateway struct {
		HttpPort       int
		WebSocketPorts cli.StringSlice
	}
	Mock struct {
		EimEndpoints cli.StringSlice
		ClientCount  int
		MessageCount int
	}
	MainDB struct {
		Driver     string
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
}
