package config

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/netutil"
)

var SystemConfig *systemConfig

func init() {
	ip, err := netutil.InternalIP()
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
		Endpoints            cli.StringSlice
		Password             string
		OfflineMessageExpire int
		OfflineDeviceExpire  int
	}
	Etcd struct {
		Endpoints cli.StringSlice
	}
	GatewaySvr struct {
		WebSocketPort int
	}
	Mock struct {
		StartUserId       int
		StartGroupId      int
		ClientCount       int
		UserMessageCount  int
		GroupMessageCount int
		SendCount         int
	}
	Database struct {
		Driver      string
		Name        string
		Connections cli.StringSlice
	}
	Mq struct {
		Endpoints cli.StringSlice
	}
	SeqSvr struct {
		RpcPort int
	}
	StorageSvr struct {
		RpcPort          int
		RegistryServices cli.StringSlice
	}
	AuthSvr struct {
		RpcPort int
		Mode    string
		OAuth2  struct {
			Endpoint     string
			ClientId     string
			ClientSecret string
		}
	}
	ApiSvr struct {
		HttpPort int
	}
	FileFlexSvr struct {
		HttpPort int
	}
	Minio struct {
		Endpoint      string
		AdminUserName string
		AdminPassword string
		UseSSL        bool
	}
}
