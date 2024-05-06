package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	MqEndpoints = "MQ_ENDPOINTS"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "mq-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:4222", "127.0.0.1:4223", "127.0.0.1:4224"),
			Usage:       "Mq地址",
			EnvVars:     []string{MqEndpoints},
			Destination: &config.SystemConfig.Mq.Endpoints,
		},
		&cli.StringSliceFlag{
			Name:        "etcd-endpoint",
			Value:       cli.NewStringSlice("127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &config.SystemConfig.Etcd.Endpoints,
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
