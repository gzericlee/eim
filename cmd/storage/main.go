package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	"eim/internal/config"
	"eim/internal/database"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/pprof"
	"eim/util/log"
	"eim/util/net"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-storage"
	app.Usage = "EIM-存储服务"
	app.Authors = []*cli.Author{
		{
			Name:  "EricLee",
			Email: "80889048@qq.com",
		},
	}
	ParseFlags(app)
	app.Action = func(c *cli.Context) error {

		//打印版本信息
		eim.Printf()

		//开启PProf服务
		pprof.EnablePProf()

		//获取随机端口
		for {
			port, err := net.GetRandomPort()
			if err != nil {
				log.Error("Error get random port", zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}
			config.SystemConfig.StorageSvr.RpcPort = port
			break
		}

		//开启Storage服务
		go func() {
			for {
				err := storagerpc.StartServer(storagerpc.Config{
					Ip:                   config.SystemConfig.LocalIp,
					Port:                 config.SystemConfig.StorageSvr.RpcPort,
					DatabaseName:         config.SystemConfig.Database.Name,
					EtcdEndpoints:        config.SystemConfig.Etcd.Endpoints.Value(),
					DatabaseDriver:       database.Driver(config.SystemConfig.Database.Driver),
					DatabaseConnection:   config.SystemConfig.Database.Connection,
					RedisEndpoints:       config.SystemConfig.Redis.Endpoints.Value(),
					RedisPassword:        config.SystemConfig.Redis.Password,
					OfflineMessageExpire: config.SystemConfig.Redis.OfflineMessageExpire,
					OfflineDeviceExpire:  config.SystemConfig.Redis.OfflineDeviceExpire,
					RegistryServices:     config.SystemConfig.StorageSvr.RegistryServices.Value(),
				})
				if err != nil {
					log.Error("Error start storage rpc server", zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.StorageSvr.RpcPort)), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.StorageSvr.RpcPort)))

		select {}

	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}

func main() {
	app := newCliApp()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v server start error: %v\n", eim.ServiceName, err)
		os.Exit(1)
	}
}
