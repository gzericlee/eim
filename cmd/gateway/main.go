package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/gateway"
	"eim/internal/gateway/websocket"
	"eim/internal/mq"
	"eim/internal/version"
	"eim/pkg/log"
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

		//打印版本信息
		version.Printf()

		//开启WS服务
		var err error
		var server *websocket.Server
		for {
			server, err = gateway.StartWebsocketServer(gateway.Config{
				Ports:          config.SystemConfig.GatewaySvr.WebSocketPorts.Value(),
				MqEndpoints:    config.SystemConfig.Mq.Endpoints.Value(),
				EtcdEndpoints:  config.SystemConfig.Etcd.Endpoints.Value(),
				RedisEndpoints: config.SystemConfig.Redis.Endpoints.Value(),
				RedisPassword:  config.SystemConfig.Redis.Password,
			})
			if err != nil {
				log.Error("WebSocket server startup error", zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}

		//初始化Nsq消费者
		for {
			consumer, err := mq.NewConsumer(config.SystemConfig.Mq.Endpoints.Value())
			if err != nil {
				goto ERROR
			}

			err = consumer.Subscribe(fmt.Sprintf(mq.MessageSendSubject, strings.Replace(config.SystemConfig.LocalIp, ".", "-", -1)), "", &websocket.SendHandler{
				Server: server,
			})
			if err != nil {
				goto ERROR
			}
			break

		ERROR:
			log.Error("Error creating mq consumers", zap.Strings("endpoints", config.SystemConfig.Mq.Endpoints.Value()), zap.Error(err))
			time.Sleep(time.Second * 5)
			continue

		}

		log.Info("Created mq consumers successfully")

		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		log.Info("PProf service started successfully", zap.String("addr", l.Addr().String()))

		log.Info(fmt.Sprintf("%v service started successfully", version.ServiceName))

		return http.Serve(l, nil)

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
