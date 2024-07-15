package main

import (
	"github.com/urfave/cli/v2"

	"github.com/gzericlee/eim/internal/config"
)

const (
	ClientCount       = "CLIENT_COUNT"
	UserMessageCount  = "USER_MESSAGE_COUNT"
	GroupMessageCount = "GROUP_MESSAGE_COUNT"

	StartUserId  = "START_USER_ID"
	StartGroupId = "START_GROUP_ID"

	SendCount = "SEND_COUNT"

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
			Name:        "start-user-id",
			Value:       1,
			Usage:       "Mock start user id",
			EnvVars:     []string{StartUserId},
			Destination: &config.SystemConfig.Mock.StartUserId,
		},
		&cli.IntFlag{
			Name:        "send-count",
			Value:       1,
			Usage:       "Mock send count",
			EnvVars:     []string{SendCount},
			Destination: &config.SystemConfig.Mock.SendCount,
		},
		&cli.IntFlag{
			Name:        "start-group-id",
			Value:       1,
			Usage:       "Mock start group id",
			EnvVars:     []string{StartGroupId},
			Destination: &config.SystemConfig.Mock.StartGroupId,
		},
		&cli.IntFlag{
			Name:        "user-message-count",
			Value:       1,
			Usage:       "Mock one client send user message count",
			EnvVars:     []string{UserMessageCount},
			Destination: &config.SystemConfig.Mock.UserMessageCount,
		},
		&cli.IntFlag{
			Name:        "group-message-count",
			Value:       1,
			Usage:       "Mock one client send group message count",
			EnvVars:     []string{GroupMessageCount},
			Destination: &config.SystemConfig.Mock.GroupMessageCount,
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
