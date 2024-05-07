package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/database"
	storagerpc "eim/internal/storage/rpc"
	"eim/internal/version"
	"eim/pkg/log"
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
		version.Printf()

		//开启Storage服务
		go func() {
			for {
				err := storagerpc.StartServer(storagerpc.Config{
					Ip:                 config.SystemConfig.LocalIp,
					Port:               config.SystemConfig.StorageSvr.RpcPort,
					DatabaseName:       config.SystemConfig.Database.Name,
					EtcdEndpoints:      config.SystemConfig.Etcd.Endpoints.Value(),
					DatabaseDriver:     database.Driver(config.SystemConfig.Database.Driver),
					DatabaseConnection: config.SystemConfig.Database.Connection,
					RedisEndpoints:     config.SystemConfig.Redis.Endpoints.Value(),
					RedisPassword:      config.SystemConfig.Redis.Password,
				})
				if err != nil {
					log.Error("Error starting rpc server", zap.Int("port", config.SystemConfig.SeqSvr.RpcPort), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", version.ServiceName))

		select {}

	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}

func main() {
	app := newCliApp()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v server startup error: %v\n", version.ServiceName, err)
		os.Exit(1)
	}
}
