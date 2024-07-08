package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"

	"eim"
	"eim/internal/config"
	"eim/internal/gateway"
	"eim/internal/gateway/handler"
	"eim/internal/gateway/server"
	"eim/internal/mq"
	"eim/pkg/exitutil"
	"eim/pkg/log"
	eimmetrics "eim/pkg/metrics"
	"eim/pkg/pprof"
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
		eim.Printf()

		//开启PProf服务
		pprof.EnablePProf()

		//开启WS服务
		wsServer, err := gateway.StartWebsocketServer(&gateway.Config{
			Ip:            config.SystemConfig.LocalIp,
			Port:          config.SystemConfig.GatewaySvr.WebSocketPort,
			MqEndpoints:   config.SystemConfig.Mq.Endpoints.Value(),
			EtcdEndpoints: config.SystemConfig.Etcd.Endpoints.Value(),
		})
		if err != nil {
			panic(fmt.Errorf("start ws server -> %w", err))
		}

		defer wsServer.Stop()

		//初始化Nsq消费者
		consumer, err := mq.NewConsumer(config.SystemConfig.Mq.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new mq consumer -> %w", err))
		}

		fmtAddr := strings.Replace(config.SystemConfig.LocalIp, ".", "-", -1)
		fmtAddr = fmt.Sprintf("%s-%d", fmtAddr, config.SystemConfig.GatewaySvr.WebSocketPort)
		err = consumer.Subscribe(fmt.Sprintf(mq.SendMessageSubjectFormat, fmtAddr), fmtAddr, &handler.SendHandler{
			Servers: []server.IServer{wsServer},
		})
		if err != nil {
			panic(fmt.Errorf("subscribe send message subject -> %w", err))
		}

		log.Info("New mq consumers successfully")

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName))

		eimmetrics.EnableMetrics(32003)

		exitutil.WaitSignal(func() {
			log.Info(fmt.Sprintf("%v service stopped successfully", eim.ServiceName))
		})

		return nil
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
