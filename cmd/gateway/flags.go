package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	HttpPort       = "HTTP_PORT"
	WebSocketPorts = "WEBSOCKET_PORTS"

	MqEndpoints = "MQ_ENDPOINTS"

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
			Destination: &config.SystemConfig.GatewaySvr.HttpPort,
		},
		&cli.StringSliceFlag{
			Name:        "websocket-ports",
			Value:       cli.NewStringSlice("10081", "10082", "10083", "10084", "10085", "10086", "10087", "10088", "10089", "10090"),
			Usage:       "Websocket ports",
			EnvVars:     []string{WebSocketPorts},
			Destination: &config.SystemConfig.GatewaySvr.WebSocketPorts,
		},
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
