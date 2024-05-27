package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	EimEndpoints = "EIM_ENDPOINTS"
	ClientCount  = "CLIENT_COUNT"
	MessageCount = "MESSAGE_COUNT"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "eim-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:10081", "127.0.0.1:10082", "127.0.0.1:10083", "127.0.0.1:10084", "127.0.0.1:10085", "127.0.0.1:10086", "127.0.0.1:10087", "127.0.0.1:10088", "127.0.0.1:10089", "127.0.0.1:10090"),
			Usage:       "EIM endpoints",
			EnvVars:     []string{EimEndpoints},
			Destination: &config.SystemConfig.Mock.EimEndpoints,
		},
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
