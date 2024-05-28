package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	EtcdEndpoints = "ETCD_ENDPOINTS"

	AuthMode = "AUTH_MODE"

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
		&cli.StringFlag{
			Name:        "auth-mode",
			Value:       "basic",
			Usage:       "Auth mode",
			EnvVars:     []string{AuthMode},
			Destination: &config.SystemConfig.AuthSvr.Mode,
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
