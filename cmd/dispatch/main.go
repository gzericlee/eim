package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim"
	"eim/internal/config"
	"eim/internal/dispatch"
	"eim/internal/mq"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/pprof"
	"eim/util/log"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-dispatch"
	app.Usage = "EIM-分发服务"
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

		//初始化Nsq消费者
		storageRpc, err := storagerpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
		if err != nil {
			panic(fmt.Errorf("new storage rpc client -> %w", err))
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

		err = consumer.Subscribe(mq.UserMessageSubject, "dispatch-user-message", dispatch.NewUserMessageHandler(storageRpc, producer))
		if err != nil {
			panic(fmt.Errorf("subscribe dispatch user message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.UserMessageSubject, "save-user-message", dispatch.NewSaveMessageHandler(storageRpc))
		if err != nil {
			panic(fmt.Errorf("subscribe save user message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.GroupMessageSubject, "dispatch-group-message", dispatch.NewGroupMessageHandler(storageRpc, producer))
		if err != nil {
			panic(fmt.Errorf("subscribe dispatch group message subject -> %w", err))
		}

		err = consumer.Subscribe(mq.GroupMessageSubject, "save-group-message", dispatch.NewSaveMessageHandler(storageRpc))
		if err != nil {
			panic(fmt.Errorf("subscribe save group message subject -> %w", err))
		}

		log.Info("New mq consumers successfully", zap.Strings("subjects", []string{mq.UserMessageSubject, mq.GroupMessageSubject}))

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName))

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
