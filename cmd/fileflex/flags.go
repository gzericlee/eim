package main

import (
	"github.com/urfave/cli/v2"

	"github.com/gzericlee/eim/internal/config"
)

const (
	HttpPort = "HTTP_PORT"

	RedisEndpoint = "REDIS_ENDPOINT"
	RedisPassword = "REDIS_PASSWORD"

	EtcdEndpoints = "ETCD_ENDPOINTS"

	MinioEndpoint      = "MINIO_ENDPOINT"
	MinioAdminUserName = "MINIO_ADMIN_USER_NAME"
	MinioAdminPassword = "MINIO_ADMIN_PASSWORD"

	LogLevel = "LOG_LEVEL"
)

func ParseFlags(app *cli.App) {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "http-port",
			Value:       10050,
			Usage:       "Http port",
			EnvVars:     []string{HttpPort},
			Destination: &config.SystemConfig.FileFlexSvr.HttpPort,
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
		&cli.StringSliceFlag{
			Name:        "etcd-endpoints",
			Value:       cli.NewStringSlice("127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"),
			Usage:       "Redis cluster endpoints",
			EnvVars:     []string{EtcdEndpoints},
			Destination: &config.SystemConfig.Etcd.Endpoints,
		},
		&cli.StringFlag{
			Name:        "minio-endpoint",
			Value:       "127.0.0.1:9000",
			Usage:       "Minio endpoint",
			EnvVars:     []string{MinioEndpoint},
			Destination: &config.SystemConfig.Minio.Endpoint,
		},
		&cli.StringFlag{
			Name:        "minio-admin-user-name",
			Value:       "minioadmin",
			Usage:       "Minio admin user name",
			EnvVars:     []string{MinioAdminUserName},
			Destination: &config.SystemConfig.Minio.AdminUserName,
		},
		&cli.StringFlag{
			Name:        "minio-admin-password",
			Value:       "minioadmin",
			Usage:       "Minio admin password",
			EnvVars:     []string{MinioAdminPassword},
			Destination: &config.SystemConfig.Minio.AdminPassword,
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
