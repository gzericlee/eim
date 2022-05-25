package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	NsqEndpoints = "NSQ_ENDPOINTS"

	MainDBDriver     = "MAIN_DB_DRIVER"
	MainDBConnection = "MAIN_DB_CONNECTION"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	LogLevel = "LOG_LEVEL"
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
		&cli.StringFlag{
			Name:        "main-db-driver",
			Value:       "mysql",
			Usage:       "Main db driver",
			EnvVars:     []string{MainDBDriver},
			Destination: &global.SystemConfig.MainDB.Driver,
		},
		&cli.StringFlag{
			Name:        "main-db-connection",
			Value:       "root:pass@word1@tcp(10.8.12.23:4000)/eim?charset=utf8mb4&parseTime=True&loc=Local",
			Usage:       "Main db connection",
			EnvVars:     []string{MainDBConnection},
			Destination: &global.SystemConfig.MainDB.Connection,
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
