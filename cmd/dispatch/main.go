package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/panjf2000/ants"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/dispatch"
	"eim/internal/mq"
	"eim/internal/redis"
	storagerpc "eim/internal/storage/rpc"
	"eim/internal/version"
	"eim/pkg/log"
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
		version.Printf()

		//初始化Nsq消费者
		for {
			storageRpc, err := storagerpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				log.Error("Error creating storage rpc client", zap.Strings("endpoints", config.SystemConfig.Etcd.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			redisManager, err := redis.NewManager(config.SystemConfig.Redis.Endpoints.Value(), config.SystemConfig.Redis.Password)
			if err != nil {
				log.Error("Error creating redis manager", zap.Strings("endpoints", config.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			producer, err := mq.NewProducer(config.SystemConfig.Nsq.Endpoints.Value())
			if err != nil {
				log.Error("Error creating mq producer", zap.Strings("endpoints", config.SystemConfig.Nsq.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			taskPool, err := ants.NewPoolPreMalloc(runtime.NumCPU() * 1000)
			if err != nil {
				log.Error("Error creating task pool", zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			consumer, err := mq.NewConsumer(config.SystemConfig.Nsq.Endpoints.Value())

			err = consumer.Subscribe(string(mq.UserMessageDispatchTopic), string(mq.MessageDispatchChannel), &dispatch.UserMessageHandler{
				StorageRpc:   storageRpc,
				RedisManager: redisManager,
				Producer:     producer,
				TaskPool:     taskPool,
			})
			if err != nil {
				goto ERROR
			}

			err = consumer.Subscribe(string(mq.GroupMessageDispatchTopic), string(mq.MessageDispatchChannel), &dispatch.GroupMessageHandler{
				StorageRpc:   storageRpc,
				RedisManager: redisManager,
				Producer:     producer,
				TaskPool:     taskPool,
			})
			if err != nil {
				goto ERROR
			}

			break

		ERROR:
			log.Error("Error creating mq consumers", zap.Error(err))
			time.Sleep(time.Second * 5)
			continue
		}

		log.Info("Created mq consumers successfully")

		log.Info(fmt.Sprintf("%v Service started successfully", version.ServiceName))

		select {}

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
