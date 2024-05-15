package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"syscall"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/mock"
	"eim/internal/version"
	"eim/util/log"
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
		version.Printf()

		log.Info("EIM GatewaySvr", zap.Strings("endpoints", config.SystemConfig.Mock.EimEndpoints.Value()))
		log.Info("Mock info", zap.Int("client_number", config.SystemConfig.Mock.ClientCount), zap.Int("send_number", config.SystemConfig.Mock.MessageCount))

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		log.Info("PProf service started successfully", zap.String("addr", l.Addr().String()))

		log.Info(fmt.Sprintf("%v service started successfully", version.ServiceName))

		mock.Do()

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

func redirectStderr(f *os.File) {
	_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}
