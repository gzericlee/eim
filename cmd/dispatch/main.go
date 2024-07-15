package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/gzericlee/eim"
	"github.com/gzericlee/eim/internal/config"
	"github.com/gzericlee/eim/internal/dispatch"
	"github.com/gzericlee/eim/internal/mq"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	"github.com/gzericlee/eim/pkg/exitutil"
	"github.com/gzericlee/eim/pkg/log"
	eimmetrics "github.com/gzericlee/eim/pkg/metrics"
	"github.com/gzericlee/eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-dispatch"
	app.Usage = "EIM-DISPATCH分发服务"
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

		bizMemberRpc, err := storagerpc.NewBizMemberClient(config.SystemConfig.Etcd.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new biz member rpc client -> %w", err))
		}

		messageRpc, err := storagerpc.NewMessageClient(config.SystemConfig.Etcd.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new message rpc client -> %w", err))
		}

		deviceRpc, err := storagerpc.NewDeviceClient(config.SystemConfig.Etcd.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new device rpc client -> %w", err))
		}

		producer, err := mq.NewProducer(config.SystemConfig.Mq.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new mq producer -> %w", err))
		}

		log.Info("New mq producers successfully")

		consumer, err := mq.NewConsumer(config.SystemConfig.Mq.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new mq consumer -> %w", err))
		}

		userMessageHandler := dispatch.NewUserMessageHandler(bizMemberRpc, messageRpc, deviceRpc, producer)

		err = consumer.Subscribe(mq.UserMessageSubject, "dispatch-user-message", userMessageHandler)
		if err != nil {
			panic(fmt.Errorf("subscribe dispatch user message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.UserMessageSubject, "save-user-message", dispatch.NewSaveMessageHandler(messageRpc))
		if err != nil {
			panic(fmt.Errorf("subscribe save user message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.GroupMessageSubject, "dispatch-group-message", dispatch.NewGroupMessageHandler(bizMemberRpc, userMessageHandler, producer))
		if err != nil {
			panic(fmt.Errorf("subscribe dispatch group message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.GroupMessageSubject, "save-group-message", dispatch.NewSaveMessageHandler(messageRpc))
		if err != nil {
			panic(fmt.Errorf("subscribe save group message subject -> %w", err))
		}

		log.Info("New mq consumers successfully", zap.Strings("subjects", []string{mq.UserMessageSubject, mq.GroupMessageSubject}))

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
