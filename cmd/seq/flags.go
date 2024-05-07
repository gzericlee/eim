package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	RpcPort = "RPC_PORT"

	DatabaseDriver     = "DATABASE_DRIVER"
	DatabaseName       = "DATABASE_NAME"
	DatabaseConnection = "DATABASE_CONNECTION"

	RedisEndpoints = "REDIS_ENDPOINTS"
	RedisPassword  = "REDIS_PASSWORD"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "rpc-port",
			Value:       10000,
			Usage:       "Rpc port",
			EnvVars:     []string{RpcPort},
			Destination: &config.SystemConfig.SeqSvr.RpcPort,
		},
		&cli.StringFlag{
			Name:        "database-driver",
			Value:       "mongodb",
			Usage:       "database driver",
			EnvVars:     []string{DatabaseDriver},
			Destination: &config.SystemConfig.Database.Driver,
		},
		&cli.StringFlag{
			Name:        "database-name",
			Value:       "eim",
			Usage:       "database name",
			EnvVars:     []string{DatabaseName},
			Destination: &config.SystemConfig.Database.Name,
		},
		&cli.StringFlag{
			Name:        "database-connection",
			Value:       "mongodb://admin:pass%40word1@127.0.0.1:27017/?authSource=admin&connect=direct",
			Usage:       "database connection",
			EnvVars:     []string{DatabaseConnection},
			Destination: &config.SystemConfig.Database.Connection,
		},
		&cli.StringSliceFlag{
			Name:        "redis-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{RedisEndpoints},
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
			Value:       "DEBUG",
			Usage:       "Log level",
			EnvVars:     []string{LogLevel},
			Destination: &config.SystemConfig.LogLevel,
		},
	}
	app.Flags = append(app.Flags, flags...)
}
