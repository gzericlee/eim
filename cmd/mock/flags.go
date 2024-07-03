package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	ClientCount  = "CLIENT_COUNT"
	MessageCount = "MESSAGE_COUNT"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "etcd-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &config.SystemConfig.Etcd.Endpoints,
		},
		&cli.IntFlag{
			Name:        "client-count",
			Value:       1,
			Usage:       "Mock client count",
			EnvVars:     []string{ClientCount},
			Destination: &config.SystemConfig.Mock.ClientCount,
		},
		&cli.IntFlag{
			Name:        "message-count",
			Value:       1,
			Usage:       "Mock one client sent message count",
			EnvVars:     []string{MessageCount},
			Destination: &config.SystemConfig.Mock.MessageCount,
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
