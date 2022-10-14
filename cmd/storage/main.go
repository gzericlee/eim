package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/build"
	"eim/internal/config"
	"eim/internal/database/maindb"
	"eim/internal/redis"
	"eim/internal/storage"
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
		build.Printf()

		//初始化日志
		log.InitLogger(log.Config{
			ConsoleEnabled: true,
			ConsoleLevel:   config.SystemConfig.LogLevel,
			ConsoleJson:    false,
			FileEnabled:    false,
			FileLevel:      config.SystemConfig.LogLevel,
			FileJson:       false,
			Directory:      "./logs/" + strings.ToLower(build.ServiceName) + "/",
			Filename:       time.Now().Format("20060102") + ".log",
			MaxSize:        200,
			MaxBackups:     10,
			MaxAge:         30,
		})

		//初始化Redis连接
		for {
			err := redis.InitRedisClusterClient(config.SystemConfig.Redis.Endpoints.Value(), config.SystemConfig.Redis.Password)
			if err != nil {
				log.Error("Error connecting to Redis cluster", zap.Strings("endpoints", config.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		log.Info("Connected Redis cluster successful")

		//初始化Tidb
		for {
			err := maindb.InitDBEngine(config.SystemConfig.MainDB.Driver, config.SystemConfig.MainDB.Connection)
			if err != nil {
				log.Error("Error connecting to Tidb", zap.String("endpoint", config.SystemConfig.MainDB.Connection), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		log.Info("Connected Tidb successful")

		//开启Storage服务
		go func() {
			err := storage.InitStorageServer(config.SystemConfig.LocalIp, config.SystemConfig.StorageSvr.RpcPort, config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				log.Error("Error starting Storage rpc server", zap.Int("port", config.SystemConfig.SeqSvr.RpcPort), zap.Error(err))
			}
		}()

		log.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

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
