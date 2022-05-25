package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"eim/build"
	"eim/global"
	"eim/internal/redis"
	"eim/internal/seq"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-seq"
	app.Usage = "EIM消息总线-ID分发服务"
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

		//初始化Redis连接
		for {
			err := redis.InitRedisClusterClient(global.SystemConfig.Redis.Endpoints.Value(), global.SystemConfig.Redis.Password)
			if err != nil {
				global.Logger.Errorf("Error connecting to Redis cluster %v : %v", global.SystemConfig.Redis.Endpoints.Value(), err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
		global.Logger.Infof("Connected Redis cluster successful")

		//开启Rpc服务
		go func() {
			err := seq.InitSeqServer(global.SystemConfig.LocalIp, global.SystemConfig.SeqSvr.RpcPort)
			if err != nil {
				global.Logger.Errorf("Seq server startup error: %s", err)
			}
		}()

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

func redirectStderr(f *os.File) {
	_ = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
}
