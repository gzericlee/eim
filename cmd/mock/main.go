package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	"eim/internal/config"
	"eim/internal/mock"
	seqrpc "eim/internal/seq/rpc"
	"eim/pkg/log"
	"eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-mock"
	app.Usage = "EIM-模拟服务"
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

		log.Info("Mock client count", zap.Int("count", config.SystemConfig.Mock.ClientCount))
		log.Info("Mock user message count", zap.Int("count", config.SystemConfig.Mock.UserMessageCount))
		log.Info("Mock group message count", zap.Int("count", config.SystemConfig.Mock.GroupMessageCount))
		log.Info("Mock start user id", zap.Int("id", config.SystemConfig.Mock.StartUserId))
		log.Info("Mock start group id", zap.Int("id", config.SystemConfig.Mock.StartGroupId))
		log.Info("Mock send count", zap.Int("count", config.SystemConfig.Mock.SendCount))

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName))

		seqRpc, err := seqrpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
		if err != nil {
			panic(err)
		}

		server := mock.NewMockServer(
			seqRpc,
			config.SystemConfig.Mock.UserMessageCount,
			config.SystemConfig.Mock.GroupMessageCount,
			config.SystemConfig.Mock.ClientCount,
			config.SystemConfig.Mock.StartUserId,
			config.SystemConfig.Mock.StartGroupId,
			config.SystemConfig.Mock.SendCount,
		)

		server.Start()

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
