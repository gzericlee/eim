package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/api"
	"eim/internal/config"
	"eim/internal/redis"
	"eim/internal/version"
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
		version.Printf()

		go func() {
			httpServer := api.HttpServer{}
			//开启Http服务
			for {
				redisManager, err := redis.NewManager(config.SystemConfig.Redis.Endpoints.Value(), config.SystemConfig.Redis.Password)
				if err != nil {
					log.Error("Error creating redis manager", zap.Strings("endpoints", config.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}

				err = httpServer.Run(api.Config{Port: config.SystemConfig.ApiSvr.HttpPort, RedisManager: redisManager})
				if err != nil {
					log.Error("ApiSvr server startup error", zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		log.Info("PProf service started successfully", zap.String("addr", l.Addr().String()))

		log.Info(fmt.Sprintf("%v Service started successfully", version.ServiceName))

		return http.Serve(l, nil)

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
