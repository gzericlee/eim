package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	EimEndpoints = "EIM_ENDPOINTS"
	ClientCount  = "CLIENT_COUNT"
	MessageCount = "MESSAGE_COUNT"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "eim-endpoints",
			Value:       cli.NewStringSlice("10.8.12.23:10081", "10.8.12.23:10082", "10.8.12.23:10083", "10.8.12.23:10084", "10.8.12.23:10085", "10.8.12.23:10086", "10.8.12.23:10087", "10.8.12.23:10088", "10.8.12.23:10089", "10.8.12.23:10090"),
			Usage:       "EIM endpoints",
			EnvVars:     []string{EimEndpoints},
			Destination: &global.SystemConfig.Mock.EimEndpoints,
		},
		&cli.IntFlag{
			Name:        "client-count",
			Value:       100000,
			Usage:       "Mock client count",
			EnvVars:     []string{ClientCount},
			Destination: &global.SystemConfig.Mock.ClientCount,
		},
		&cli.IntFlag{
			Name:        "message-count",
			Value:       5,
			Usage:       "Mock one client sent message count",
			EnvVars:     []string{MessageCount},
			Destination: &global.SystemConfig.Mock.MessageCount,
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
