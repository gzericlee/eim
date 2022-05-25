package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"eim/build"
	"eim/global"
	"eim/internal/database/maindb"
	"eim/internal/nsq/consumer"
	"eim/internal/redis"
	"eim/model"
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

		//打印编译信息
		build.Printf()

		//初始化日志
		global.InitLogger()

		//初始化Redis连接
		for {
			err := redis.InitRedisClusterClient(global.SystemConfig.Redis.Endpoints.Value(), global.SystemConfig.Redis.Password)
			if err != nil {
				global.Logger.Errorf("Error connecting to Redis cluster %v : %v", global.SystemConfig.Redis.Endpoints.Value(), err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Connected Redis cluster successful")

		//初始化Tidb
		for {
			err := maindb.InitDBEngine(global.SystemConfig.MainDB.Driver, global.SystemConfig.MainDB.Connection)
			if err != nil {
				global.Logger.Errorf("Error connecting to Tidb %v : %v", global.SystemConfig.MainDB.Connection, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Connected Tidb successful")

		//初始化Nsq消费者
		for {
			err := consumer.InitConsumers(map[string][]string{
				model.DeviceStoreTopic:  []string{model.DeviceStoreChannel},
				model.MessageStoreTopic: []string{model.MessageStoreChannel},
			}, global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Errorf("Error createing Nsq %v consumers: %v", global.SystemConfig.Nsq.Endpoints.Value(), err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Created Nsq consumers successful")

		global.Logger.Infof("%v service started successful", build.ServiceName)

		select {}

	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}

func main() {
	app := newCliApp()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v server startup error: %v\n", build.ServiceName, err)
		os.Exit(1)
	}
}
