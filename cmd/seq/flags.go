package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	RpcPort = "RPC_PORT"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "rpc-port",
			Value:       10000,
			Usage:       "Rpc port",
			EnvVars:     []string{RpcPort},
			Destination: &global.SystemConfig.SeqSvr.RpcPort,
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
			Name:        "log-level",
			Value:       "INFO",
			Usage:       "Log level",
			EnvVars:     []string{LogLevel},
			Destination: &global.SystemConfig.LogLevel,
		},
	}
	app.Flags = append(app.Flags, flags...)
}
