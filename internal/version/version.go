package version

import (
	"strings"
	"time"

	"github.com/dimiro1/banner"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"

	"eim/internal/config"
	"eim/util/log"
)

var (
	ServiceName string
	Branch      string
	Commit      string
	Date        string
)

func initLogger() {
	//初始化日志
	log.SetConfig(log.Config{
		ConsoleEnabled: true,
		ConsoleLevel:   config.SystemConfig.LogLevel,
		ConsoleJson:    false,
		FileEnabled:    false,
		FileLevel:      config.SystemConfig.LogLevel,
		FileJson:       false,
		Directory:      "./logs/" + strings.ToLower(ServiceName) + "/",
		Filename:       time.Now().Format("20060102") + ".log",
		MaxSize:        200,
		MaxBackups:     10,
		MaxAge:         30,
	})
}

func Printf() {
	if ServiceName == "" {
		ServiceName = "EIM-?"
	}
	if Branch == "" {
		Branch = "master"
	}
	if Commit == "" {
		Commit = "dev"
	}
	if Date == "" {
		Date = time.Now().Format("2006-01-02 15:04:05")
	}

	initLogger()

	template := `
{{ .Title "` + ServiceName + `" "" 0 }}
{{ .AnsiColor.BrightCyan }}OpenSource Instant Messaging Server{{ .AnsiColor.Default }}
Branch: ` + Branch + `
Commit: ` + Commit + `
Date:   ` + strings.Replace(Date, "T", " ", -1) + `

`
	banner.InitString(colorable.NewColorableStdout(), true, true, template)
}
