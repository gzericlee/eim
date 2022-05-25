package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"eim/build"
	"eim/global"
	"eim/internal/mock"

	"github.com/urfave/cli/v2"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-mock"
	app.Usage = "EIM消息总线-模拟服务"
	app.Authors = []*cli.Author{
		{
			Name:  "LiRui",
			Email: "lirui@gz-mstc.com",
		},
	}
	ParseFlags(app)
	app.Action = func(c *cli.Context) error {

		//打印编译信息
		build.Printf()

		//初始化日志
		global.InitLogger()

		global.Logger.Infof("Mock service endpoints: %v", global.SystemConfig.Mock.EimEndpoints.Value())
		global.Logger.Infof("Mock client count: %v", global.SystemConfig.Mock.ClientCount)
		global.Logger.Infof("Mock one client sent message count: %v", global.SystemConfig.Mock.MessageCount)

		global.Logger.Infof("%v service started successful", build.ServiceName)

		mock.Do()

		select {}

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

func redirectStderr(f *os.File) {
	_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}
