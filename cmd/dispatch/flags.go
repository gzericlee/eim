package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	NsqEndpoints = "NSQ_ENDPOINTS"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "nsq-endpoints",
			Value:       cli.NewStringSlice("10.8.12.23:4161", "10.8.12.23:4261"),
			Usage:       "Nsqlookupd地址",
			EnvVars:     []string{NsqEndpoints},
			Destination: &global.SystemConfig.Nsq.Endpoints,
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
