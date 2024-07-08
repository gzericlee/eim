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
	"eim/pkg/exitutil"
	"eim/pkg/log"
	"eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-file-flex"
	app.Usage = "EIM-FILE-FLEX服务"
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

		httpServer := api.HttpServer{}

		go func() {
			//开启Http服务
			redisManager, err := redis.NewManager(redis.Config{
				RedisEndpoints: config.SystemConfig.Redis.Endpoints.Value(),
				RedisPassword:  config.SystemConfig.Redis.Password,
			})
			if err != nil {
				panic(fmt.Errorf("new redis manager -> %w", err))
			}

			log.Info("New redis manager successfully")

			_ = httpServer.Run(api.Config{
				Port:         config.SystemConfig.ApiSvr.HttpPort,
				RedisManager: redisManager})
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.Int("port", config.SystemConfig.FileFlexSvr.HttpPort))

		exitutil.WaitSignal(func() {
			httpServer.Stop()
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
