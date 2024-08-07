package main

import (
	"github.com/urfave/cli/v2"

	"github.com/gzericlee/eim/internal/config"
)

const (
	EtcdEndpoints = "ETCD_ENDPOINTS"

	AuthMode = "AUTH_MODE"

	OAuth2ClientId     = "OAUTH2_CLIENT_ID"
	OAuth2ClientSecret = "OAUTH2_CLIENT_SECRET"
	OAuth2Endpoint     = "OAUTH2_ENDPOINT"

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
			Name:        "oauth2-endpoint",
			Value:       "",
			Usage:       "OAuth2 endpoint",
			EnvVars:     []string{OAuth2Endpoint},
			Destination: &config.SystemConfig.AuthSvr.OAuth2.Endpoint,
		},
		&cli.StringFlag{
			Name:        "oauth2-client-id",
			Value:       "client_id",
			Usage:       "OAuth2 client id",
			EnvVars:     []string{OAuth2ClientId},
			Destination: &config.SystemConfig.AuthSvr.OAuth2.ClientId,
		},
		&cli.StringFlag{
			Name:        "oauth2-client-secret",
			Value:       "client_secret",
			Usage:       "OAuth2 client secret",
			EnvVars:     []string{OAuth2ClientSecret},
			Destination: &config.SystemConfig.AuthSvr.OAuth2.ClientSecret,
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
