package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	"eim/internal/config"
	seqrpc "eim/internal/seq/rpc"
	"eim/pkg/pprof"
	"eim/util/log"
	"eim/util/net"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-seq"
	app.Usage = "EIM-ID分发服务"
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
		for {
			port, err := net.GetRandomPort()
			if err != nil {
				log.Error("Error get random port", zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}
			config.SystemConfig.SeqSvr.RpcPort = port
			break
		}

		//开启Rpc服务
		go func() {
			for {
				err := seqrpc.StartServer(seqrpc.Config{
					Ip:             config.SystemConfig.LocalIp,
					Port:           config.SystemConfig.SeqSvr.RpcPort,
					EtcdEndpoints:  config.SystemConfig.Etcd.Endpoints.Value(),
					RedisEndpoints: config.SystemConfig.Redis.Endpoints.Value(),
					RedisPassword:  config.SystemConfig.Redis.Password,
				})
				if err != nil {
					log.Error("Error start seq rpc server", zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.SeqSvr.RpcPort)), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.String("addr", fmt.Sprintf("%v:%v", config.SystemConfig.LocalIp, config.SystemConfig.SeqSvr.RpcPort)))

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

func redirectStderr(f *os.File) {
	_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}
