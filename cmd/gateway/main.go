package main

import (
	"fmt"
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
	"eim/pkg/pprof"
	"eim/util/log"
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

		//开启PProf服务
		pprof.EnablePProf()

		//开启WS服务
		var err error
		var server *websocket.Server
		for {
			server, err = gateway.StartWebsocketServer(gateway.Config{
				Ports:         config.SystemConfig.GatewaySvr.WebSocketPorts.Value(),
				MqEndpoints:   config.SystemConfig.Mq.Endpoints.Value(),
				EtcdEndpoints: config.SystemConfig.Etcd.Endpoints.Value(),
			})
			if err != nil {
				log.Error("Error start webSocket server", zap.Error(err))
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
			log.Error("Error new mq consumers", zap.Strings("endpoints", config.SystemConfig.Mq.Endpoints.Value()), zap.Error(err))
			time.Sleep(time.Second * 5)
			continue

		}

		log.Info("New mq consumers successfully")

		log.Info(fmt.Sprintf("%v service started successfully", version.ServiceName))

		select {}

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
