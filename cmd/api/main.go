package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/api"
	"eim/internal/build"
	"eim/internal/config"
	"eim/internal/redis"
	"eim/pkg/log"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-api"
	app.Usage = "EIM-API服务"
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

		go func() {
			httpServer := api.HttpServer{}
			//开启Http服务
			for {
				err := httpServer.Run(api.Config{Port: config.SystemConfig.ApiSvr.HttpPort})
				if err != nil {
					log.Error("ApiSvr server startup error", zap.Error(err))
					time.Sleep(time.Second)
					continue
				}
				break
			}
		}()

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		log.Info("PProf service started successful", zap.String("addr", l.Addr().String()))

		log.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

		return http.Serve(l, nil)

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
