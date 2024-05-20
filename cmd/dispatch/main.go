package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/dispatch"
	"eim/internal/mq"
	"eim/internal/redis"
	storagerpc "eim/internal/storage/rpc"
	"eim/internal/version"
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
		version.Printf()

		//初始化Nsq消费者
		for {
			storageRpc, err := storagerpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				log.Error("Error new storage rpc client", zap.Strings("endpoints", config.SystemConfig.Etcd.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			redisManager, err := redis.NewManager(redis.Config{
				RedisEndpoints:       config.SystemConfig.Redis.Endpoints.Value(),
				RedisPassword:        config.SystemConfig.Redis.Password,
				OfflineMessageExpire: time.Hour * 24 * time.Duration(config.SystemConfig.Redis.OfflineMessageExpire),
				OfflineDeviceExpire:  time.Hour * 24 * time.Duration(config.SystemConfig.Redis.OfflineDeviceExpire),
			})
			if err != nil {
				log.Error("Error new redis manager", zap.Strings("endpoints", config.SystemConfig.Redis.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			producer, err := mq.NewProducer(config.SystemConfig.Mq.Endpoints.Value())
			if err != nil {
				log.Error("Error new mq producer", zap.Strings("endpoints", config.SystemConfig.Mq.Endpoints.Value()), zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			log.Info("New mq producers successfully")

			consumer, err := mq.NewConsumer(config.SystemConfig.Mq.Endpoints.Value())

			err = consumer.Subscribe(mq.MessageDispatchSubject, mq.DispatchQueue, &dispatch.UserMessageHandler{
				StorageRpc:   storageRpc,
				RedisManager: redisManager,
				Producer:     producer,
			})
			if err != nil {
				log.Error("Error new mq consumers", zap.Error(err))
				time.Sleep(time.Second * 5)
				continue
			}

			log.Info("New mq consumers successfully")

			break
		}

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
