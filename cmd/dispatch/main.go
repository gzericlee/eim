package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"

	"eim/build"
	"eim/global"
	"eim/internal/nsq/consumer"
	"eim/internal/nsq/producer"
	"eim/model"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-dispatch"
	app.Usage = "EIM消息总线-分发服务"
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

		//初始化Nsq生产者
		for {
			err := producer.InitProducers(global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Errorf("Error createing Nsq %v producers: %v", global.SystemConfig.Nsq.Endpoints.Value(), err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Created Nsq producers successful")

		//初始化Nsq消费者
		for {
			err := consumer.InitConsumers(map[string][]string{
				model.MessageDispatchTopic: []string{model.MessageDispatchChannel},
			}, global.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				global.Logger.Errorf("Error createing Nsq consumers %v", err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Created Nsq consumers successful")

		global.Logger.Infof("%v service started successful", build.ServiceName)

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
