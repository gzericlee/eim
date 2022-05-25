package global

import (
	"strings"
	"time"

	"eim/build"
	"eim/pkg/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.Configure(log.Config{
		ConsoleEnabled: true,
		ConsoleLevel:   SystemConfig.LogLevel,
		ConsoleJson:    false,
		FileEnabled:    true,
		FileLevel:      "DEBUG",
		FileJson:       false,
		Directory:      "./logs/" + strings.ToLower(build.ServiceName) + "/",
		Filename:       time.Now().Format("20060102") + ".log",
		MaxSize:        200,
		MaxBackups:     10,
		MaxAge:         30,
	})
}
