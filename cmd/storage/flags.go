package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	RpcPort = "RPC_PORT"

	NsqEndpoints = "NSQ_ENDPOINTS"

	MainDBDriver     = "MAIN_DB_DRIVER"
	MainDBConnection = "MAIN_DB_CONNECTION"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "rpc-port",
			Value:       11000,
			Usage:       "Rpc port",
			EnvVars:     []string{RpcPort},
			Destination: &config.SystemConfig.StorageSvr.RpcPort,
		},
		&cli.StringSliceFlag{
			Name:        "nsq-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:4161", "127.0.0.1:4261"),
			Usage:       "Nsqlookupd endpoints",
			EnvVars:     []string{NsqEndpoints},
			Destination: &config.SystemConfig.Nsq.Endpoints,
		},
		&cli.StringFlag{
			Name:        "main-db-driver",
			Value:       "mysql",
			Usage:       "Main db driver",
			EnvVars:     []string{MainDBDriver},
			Destination: &config.SystemConfig.MainDB.Driver,
		},
		&cli.StringFlag{
			Name:        "main-db-connection",
			Value:       "root:pass@word1@tcp(127.0.0.1:4000)/eim?charset=utf8mb4&parseTime=True&loc=Local",
			Usage:       "Main db connection",
			EnvVars:     []string{MainDBConnection},
			Destination: &config.SystemConfig.MainDB.Connection,
		},
		&cli.StringSliceFlag{
			Name:        "redis-endpoint",
			Value:       cli.NewStringSlice("127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{RedisEndpoint},
			Destination: &config.SystemConfig.Redis.Endpoints,
		},
		&cli.StringFlag{
			Name:        "redis-password",
			Value:       "pass@word1",
			Usage:       "Redis passwd",
			EnvVars:     []string{RedisPassword},
			Destination: &config.SystemConfig.Redis.Password,
		},
		&cli.StringSliceFlag{
			Name:        "etcd-endpoint",
			Value:       cli.NewStringSlice("127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &config.SystemConfig.Etcd.Endpoints,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "INFO",
			Usage:       "Log level",
			EnvVars:     []string{LogLevel},
			Destination: &config.SystemConfig.LogLevel,
		},
	}
	app.Flags = append(app.Flags, flags...)
}
