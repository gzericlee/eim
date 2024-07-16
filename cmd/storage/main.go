package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/gzericlee/eim"
	"github.com/gzericlee/eim/internal/config"
	"github.com/gzericlee/eim/internal/database"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	"github.com/gzericlee/eim/pkg/exitutil"
	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/netutil"
	"github.com/gzericlee/eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-storage"
	app.Usage = "eim-存储服务"
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
		port, err := netutil.RandomPort()
		if err != nil {
			panic(fmt.Errorf("get random port -> %w", err))
		}
		config.SystemConfig.StorageSvr.RpcPort = port

		//开启Storage服务
		go func() {
			err := storagerpc.StartServer(&storagerpc.Config{
				Ip:                   config.SystemConfig.LocalIp,
				Port:                 config.SystemConfig.StorageSvr.RpcPort,
				DatabaseName:         config.SystemConfig.Database.Name,
				EtcdEndpoints:        config.SystemConfig.Etcd.Endpoints.Value(),
				DatabaseDriver:       database.Driver(config.SystemConfig.Database.Driver),
				DatabaseConnections:  config.SystemConfig.Database.Connections.Value(),
				RedisEndpoints:       config.SystemConfig.Redis.Endpoints.Value(),
				RedisPassword:        config.SystemConfig.Redis.Password,
				OfflineMessageExpire: config.SystemConfig.Redis.OfflineMessageExpire,
				OfflineDeviceExpire:  config.SystemConfig.Redis.OfflineDeviceExpire,
				RegistryServices:     config.SystemConfig.StorageSvr.RegistryServices.Value(),
			})
			if err != nil {
				panic(fmt.Errorf("start storage rpc server -> %w", err))
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.StorageSvr.RpcPort)))

		exitutil.WaitSignal(func() {
			log.Info(fmt.Sprintf("%v service stopped successfully", eim.ServiceName))
		})

		return nil
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
