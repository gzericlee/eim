package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/build"
	"eim/global"
	"eim/internal/nsq/consumer"
	"eim/internal/redis"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-auth"
	app.Usage = "EIM-鉴权服务"
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
				global.Logger.Error("Error connecting to Redis cluster", zap.Strings("endpoints", global.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Info("Connected Redis cluster successful")

		//初始化Nsq消费者
		for {
			err := consumer.InitConsumers(map[string][]string{}, global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Error("Error creating Nsq consumers", zap.Strings("endpoints", global.SystemConfig.Nsq.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Info("Created Nsq consumers successful")

		global.Logger.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

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
