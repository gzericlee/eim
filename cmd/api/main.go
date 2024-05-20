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
	"eim/util/log"
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
				redisManager, err := redis.NewManager(redis.Config{
					RedisEndpoints: config.SystemConfig.Redis.Endpoints.Value(),
					RedisPassword:  config.SystemConfig.Redis.Password,
				})
				if err != nil {
					log.Error("Error new redis manager", zap.Strings("endpoints", config.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}

				log.Info("New redis manager successfully")

				err = httpServer.Run(api.Config{Port: config.SystemConfig.ApiSvr.HttpPort, RedisManager: redisManager})
				if err != nil {
					log.Error("Error start api server", zap.Error(err))
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

		log.Info(fmt.Sprintf("%v service started successfully", version.ServiceName), zap.Int("port", config.SystemConfig.ApiSvr.HttpPort))

		return http.Serve(l, nil)

	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}

func main() {
	app := newCliApp()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v server start error: %v\n", version.ServiceName, err)
		os.Exit(1)
	}
}
