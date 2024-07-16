package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/gzericlee/eim"
	"github.com/gzericlee/eim/internal/auth"
	authrpc "github.com/gzericlee/eim/internal/auth/rpc"
	"github.com/gzericlee/eim/internal/config"
	"github.com/gzericlee/eim/pkg/exitutil"
	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/netutil"
	oauth2lib "github.com/gzericlee/eim/pkg/oauth2"
	"github.com/gzericlee/eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-auth"
	app.Usage = "EIM-AUTH鉴权服务"
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
		port, err := netutil.RandomPort()
		if err != nil {
			panic(fmt.Errorf("get random port -> %w", err))
		}
		config.SystemConfig.AuthSvr.RpcPort = port

		//开启Rpc服务
		go func() {
			var oauth2Client oauth2lib.Client
			if auth.Mode(config.SystemConfig.AuthSvr.Mode) == auth.OAuth2Mode {
				oauth2Client, err = oauth2lib.NewClient(&oauth2lib.Config{
					Version:      oauth2lib.V5,
					Endpoint:     config.SystemConfig.AuthSvr.OAuth2.Endpoint,
					ClientId:     config.SystemConfig.AuthSvr.OAuth2.ClientId,
					ClientSecret: config.SystemConfig.AuthSvr.OAuth2.ClientSecret,
				})
				if err != nil {
					panic(fmt.Errorf("new oauth2 client -> %w", err))
				}
			}

			err = authrpc.StartServer(&authrpc.Config{
				Ip:            config.SystemConfig.LocalIp,
				Port:          config.SystemConfig.AuthSvr.RpcPort,
				EtcdEndpoints: config.SystemConfig.Etcd.Endpoints.Value(),
				AuthMode:      auth.Mode(config.SystemConfig.AuthSvr.Mode),
				Oauth2Client:  oauth2Client,
			})
			if err != nil {
				panic(fmt.Errorf("start auth rpc server -> %w", err))
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.AuthSvr.RpcPort)))

		exitutil.WaitSignal(func() {
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
