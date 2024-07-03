package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	authrpc "eim/internal/auth/rpc"
	"eim/internal/config"
	"eim/pkg/pprof"
	"eim/util/log"
	"eim/util/net"
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

		//打印版本信息
		eim.Printf()

		//开启PProf服务
		pprof.EnablePProf()

		//获取随机端口
		port, err := net.RandomPort()
		if err != nil {
			panic(fmt.Errorf("get random port -> %w", err))
		}
		config.SystemConfig.AuthSvr.RpcPort = port

		//开启Rpc服务
		go func() {
			err := authrpc.StartServer(authrpc.Config{
				Ip:            config.SystemConfig.LocalIp,
				Port:          config.SystemConfig.AuthSvr.RpcPort,
				EtcdEndpoints: config.SystemConfig.Etcd.Endpoints.Value(),
			})
			if err != nil {
				panic(fmt.Errorf("start auth rpc server -> %w", err))
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.AuthSvr.RpcPort)))

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
