package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	NsqEndpoints = "NSQ_ENDPOINTS"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "nsq-endpoints",
			Value:       cli.NewStringSlice("10.8.12.23:4161", "10.8.12.23:4261"),
			Usage:       "Nsqlookupd endpoints",
			EnvVars:     []string{NsqEndpoints},
			Destination: &global.SystemConfig.Nsq.Endpoints,
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
	}
	app.Flags = append(app.Flags, flags...)
}
