package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	HttpPort       = "HTTP_PORT"
	WebSocketPorts = "WEBSOCKET_PORTS"

	NsqEndpoints = "NSQ_ENDPOINTS"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "http-port",
			Value:       10081,
			Usage:       "Http port",
			EnvVars:     []string{HttpPort},
			Destination: &global.SystemConfig.Gateway.HttpPort,
		},
		&cli.StringSliceFlag{
			Name:        "websocket-ports",
			Value:       cli.NewStringSlice("10081", "10082", "10083", "10084", "10085", "10086", "10087", "10088", "10089", "10090"),
			Usage:       "Websocket ports",
			EnvVars:     []string{WebSocketPorts},
			Destination: &global.SystemConfig.Gateway.WebSocketPorts,
		},
		&cli.StringSliceFlag{
			Name:        "nsq-endpoints",
			Value:       cli.NewStringSlice("10.8.12.23:4161", "10.8.12.23:4261"),
			Usage:       "Nsqlookupd endpoints",
			EnvVars:     []string{NsqEndpoints},
			Destination: &global.SystemConfig.Nsq.Endpoints,
		},
		&cli.StringSliceFlag{
			Name:        "etcd-endpoint",
			Value:       cli.NewStringSlice("10.8.12.23:2379", "10.8.12.23:2479", "10.8.12.23:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &global.SystemConfig.Etcd.Endpoints,
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
