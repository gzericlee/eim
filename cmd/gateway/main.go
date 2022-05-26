package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/build"
	"eim/global"
	"eim/internal/gateway/websocket"
	"eim/internal/nsq/consumer"
	"eim/internal/nsq/producer"
	"eim/model"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-gateway"
	app.Usage = "EIM-网关服务"
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

		//开启WS服务
		err := websocket.InitGatewayServer("0.0.0.0", global.SystemConfig.Gateway.WebSocketPorts.Value())
		if err != nil {
			global.Logger.Error("Gateway server startup error", zap.Error(err))
		}

		//初始化Nsq生产者
		for {
			err := producer.InitProducers(global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Error("Error creating Nsq %v producers", zap.Strings("endpoints", global.SystemConfig.Nsq.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Info("Created Nsq producers successful")

		//初始化Nsq消费者
		for {
			err := consumer.InitConsumers(map[string][]string{
				model.MessageSendTopic: []string{global.SystemConfig.LocalIp},
			}, global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Error("Error creating Nsq consumers", zap.Strings("endpoints", global.SystemConfig.Nsq.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Info("Created Nsq consumers successful")

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		global.Logger.Info("PProf service started successful", zap.String("addr", l.Addr().String()))

		global.Logger.Info(fmt.Sprintf("%v Service started successful", build.ServiceName))

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
