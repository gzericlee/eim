package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"syscall"

	"go.uber.org/zap"

	"eim/build"
	"eim/global"
	"eim/internal/mock"

	"github.com/urfave/cli/v2"
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

		//打印编译信息
		build.Printf()

		//初始化日志
		global.InitLogger()

		global.Logger.Info("EIM Gateway", zap.Strings("endpoints", global.SystemConfig.Mock.EimEndpoints.Value()))
		global.Logger.Info("Mock info", zap.Int("client_number", global.SystemConfig.Mock.ClientCount), zap.Int("send_number", global.SystemConfig.Mock.MessageCount))

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		global.Logger.Info("PProf service started successful", zap.String("addr", l.Addr().String()))

		global.Logger.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

		mock.Do()

		return http.Serve(l, nil)
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
