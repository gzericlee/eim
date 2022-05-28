package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	RpcPort = "RPC_PORT"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	AuthMode = "AUTH_MODE"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "rpc-port",
			Value:       12000,
			Usage:       "Rpc port",
			EnvVars:     []string{RpcPort},
			Destination: &global.SystemConfig.AuthSvr.RpcPort,
		},
		&cli.StringSliceFlag{
			Name:        "etcd-endpoint",
			Value:       cli.NewStringSlice("10.8.12.23:2379", "10.8.12.23:2479", "10.8.12.23:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &global.SystemConfig.Etcd.Endpoints,
		},
		&cli.StringSliceFlag{
			Name:        "redis-endpoint",
			Value:       cli.NewStringSlice("10.8.12.23:7001", "10.8.12.23:7002", "10.8.12.23:7003"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{RedisEndpoint},
			Destination: &global.SystemConfig.Redis.Endpoints,
		},
		&cli.StringFlag{
			Name:        "redis-password",
			Value:       "pass@word1",
			Usage:       "Redis passwd",
			EnvVars:     []string{RedisPassword},
			Destination: &global.SystemConfig.Redis.Password,
		},
		&cli.StringFlag{
			Name:        "auth-mode",
			Value:       "basic",
			Usage:       "Auth mode",
			EnvVars:     []string{AuthMode},
			Destination: &global.SystemConfig.AuthSvr.Mode,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "INFO",
			Usage:       "Log level",
			EnvVars:     []string{LogLevel},
			Destination: &global.SystemConfig.LogLevel,
		},
	}
	app.Flags = append(app.Flags, flags...)
}
