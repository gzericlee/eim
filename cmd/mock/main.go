package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/build"
	"eim/internal/config"
	"eim/internal/mock"
	"eim/pkg/log"
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
		build.Printf()

		//初始化日志
		log.InitLogger(log.Config{
			ConsoleEnabled: true,
			ConsoleLevel:   config.SystemConfig.LogLevel,
			ConsoleJson:    false,
			FileEnabled:    false,
			FileLevel:      config.SystemConfig.LogLevel,
			FileJson:       false,
			Directory:      "./logs/" + strings.ToLower(build.ServiceName) + "/",
			Filename:       time.Now().Format("20060102") + ".log",
			MaxSize:        200,
			MaxBackups:     10,
			MaxAge:         30,
		})

		log.Info("EIM GatewaySvr", zap.Strings("endpoints", config.SystemConfig.Mock.EimEndpoints.Value()))
		log.Info("Mock info", zap.Int("client_number", config.SystemConfig.Mock.ClientCount), zap.Int("send_number", config.SystemConfig.Mock.MessageCount))

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		log.Info("PProf service started successful", zap.String("addr", l.Addr().String()))

		log.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

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
