package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	"eim/internal/api"
	"eim/internal/config"
	"eim/internal/redis"
	"eim/pkg/pprof"
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
		eim.Printf()

		//开启PProf服务
		pprof.EnablePProf()

		go func() {
			httpServer := api.HttpServer{}
			//开启Http服务
			redisManager, err := redis.NewManager(redis.Config{
				RedisEndpoints: config.SystemConfig.Redis.Endpoints.Value(),
				RedisPassword:  config.SystemConfig.Redis.Password,
			})
			if err != nil {
				panic(fmt.Errorf("new redis manager -> %w", err))
			}

			log.Info("New redis manager successfully")

			err = httpServer.Run(api.Config{Port: config.SystemConfig.ApiSvr.HttpPort, RedisManager: redisManager})
			if err != nil {
				panic(fmt.Errorf("run http server -> %w", err))
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.Int("port", config.SystemConfig.ApiSvr.HttpPort))

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
