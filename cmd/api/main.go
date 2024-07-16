package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/gzericlee/eim"
	"github.com/gzericlee/eim/internal/api"
	authrpc "github.com/gzericlee/eim/internal/auth/rpc"
	"github.com/gzericlee/eim/internal/config"
	"github.com/gzericlee/eim/internal/minio"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc"
	"github.com/gzericlee/eim/pkg/exitutil"
	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/pprof"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "eim-api"
	app.Usage = "EIM-API服务"
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

		httpServer := api.HttpServer{}

		go func() {
			tenantRpc, err := storagerpc.NewTenantClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				panic(fmt.Errorf("new tenant rpc client -> %w", err))
			}

			gatewayRpc, err := storagerpc.NewGatewayClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				panic(fmt.Errorf("new gateway rpc client -> %w", err))
			}

			authRpc, err := authrpc.NewAuthClient(config.SystemConfig.Etcd.Endpoints.Value())
			if err != nil {
				panic(fmt.Errorf("new auth rpc client -> %w", err))
			}

			minioManager, err := minio.NewManager(&minio.Config{
				Endpoint:        config.SystemConfig.Minio.Endpoint,
				AccessKeyId:     config.SystemConfig.Minio.AdminUserName,
				SecretAccessKey: config.SystemConfig.Minio.AdminPassword,
				UseSSL:          config.SystemConfig.Minio.UseSSL,
			})
			if err != nil {
				panic(fmt.Errorf("new minio manager -> %w", err))
			}

			log.Info("New redis manager successfully")

			_ = httpServer.Run(&api.Config{
				Port:         config.SystemConfig.ApiSvr.HttpPort,
				TenantRpc:    tenantRpc,
				GatewayRpc:   gatewayRpc,
				AuthRpc:      authRpc,
				MinioManager: minioManager,
			})
		}()

		log.Info(fmt.Sprintf("%v service started successfully", eim.ServiceName), zap.Int("port", config.SystemConfig.ApiSvr.HttpPort))

		exitutil.WaitSignal(func() {
			httpServer.Stop()
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
