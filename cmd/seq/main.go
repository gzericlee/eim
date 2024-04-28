package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	seqrpc "eim/internal/seq/rpc"
	"eim/internal/version"
	"eim/pkg/log"
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
		version.Printf()

		//开启Rpc服务
		go func() {
			for {
				err := seqrpc.StartServer(seqrpc.Config{
					Ip:            config.SystemConfig.LocalIp,
					Port:          config.SystemConfig.SeqSvr.RpcPort,
					EtcdEndpoints: config.SystemConfig.Etcd.Endpoints.Value(),
				})
				if err != nil {
					log.Error("Error starting rpc server", zap.Int("port", config.SystemConfig.SeqSvr.RpcPort), zap.Error(err))
					time.Sleep(time.Second * 5)
					continue
				}
				break
			}
		}()

		log.Info(fmt.Sprintf("%v Service started successfully", version.ServiceName))

		select {}

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

func redirectStderr(f *os.File) {
	_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}
