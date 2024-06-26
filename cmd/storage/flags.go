package main

import (
	"github.com/urfave/cli/v2"

	"eim/internal/config"
)

const (
	MqEndpoints = "MQ_ENDPOINTS"

	DatabaseDriver     = "DATABASE_DRIVER"
	DatabaseName       = "DATABASE_NAME"
	DatabaseConnection = "DATABASE_CONNECTION"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	OfflineDeviceExpire  = "OFFLINE_DEVICE_EXPIRE"
	OfflineMessageExpire = "OFFLINE_MESSAGE_EXPIRE"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	RegistryServices = "REGISTRY_SERVICES"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "mq-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:4222", "127.0.0.1:4223", "127.0.0.1:4224"),
			Usage:       "Mq地址",
			EnvVars:     []string{MqEndpoints},
			Destination: &config.SystemConfig.Mq.Endpoints,
		},
		&cli.StringFlag{
			Name:        "database-driver",
			Value:       "mongodb",
			Usage:       "database driver",
			EnvVars:     []string{DatabaseDriver},
			Destination: &config.SystemConfig.Database.Driver,
		},
		&cli.StringFlag{
			Name:        "database-name",
			Value:       "eim",
			Usage:       "database name",
			EnvVars:     []string{DatabaseName},
			Destination: &config.SystemConfig.Database.Name,
		},
		&cli.StringFlag{
			Name:        "database-connection",
			Value:       "mongodb://admin:pass%40word1@127.0.0.1:27017/?authSource=admin&connect=direct",
			Usage:       "database connection",
			EnvVars:     []string{DatabaseConnection},
			Destination: &config.SystemConfig.Database.Connection,
		},
		&cli.StringSliceFlag{
			Name:        "redis-endpoint",
			Value:       cli.NewStringSlice("127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004"),
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
		&cli.IntFlag{
			Name:        "offline-message-expire",
			Value:       30,
			Usage:       "Offline message expire",
			EnvVars:     []string{OfflineMessageExpire},
			Destination: &config.SystemConfig.Redis.OfflineMessageExpire,
		},
		&cli.IntFlag{
			Name:        "offline-device-expire",
			Value:       30,
			Usage:       "Offline device expire",
			EnvVars:     []string{OfflineDeviceExpire},
			Destination: &config.SystemConfig.Redis.OfflineDeviceExpire,
		},
		&cli.StringSliceFlag{
			Name:        "etcd-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &config.SystemConfig.Etcd.Endpoints,
		},
		&cli.StringSliceFlag{
			Name:        "registry-services",
			Value:       cli.NewStringSlice("device", "message", "biz", "biz_member", "gateway", "segment"),
			Usage:       "Registry services",
			EnvVars:     []string{RegistryServices},
			Destination: &config.SystemConfig.StorageSvr.RegistryServices,
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
