package main

import (
	"github.com/urfave/cli/v2"

	"eim/global"
)

const (
	HttpPort       = "HTTP_PORT"
	WebSocketPorts = "WEBSOCKET_PORTS"

	NsqEndpoints = "NSQ_ENDPOINTS"

	SeqEndpoint = "SEQ_ENDPOINT"

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
		&cli.StringFlag{
			Name:        "seq-endpoint",
			Value:       "10.8.12.23:10000",
			Usage:       "Seq server endpoint",
			EnvVars:     []string{SeqEndpoint},
			Destination: &global.SystemConfig.SeqSvr.Endpoint,
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
