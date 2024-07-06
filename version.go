package eim

import (
	"strings"
	"time"

	"github.com/dimiro1/banner"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"

	"eim/pkg/log"

	"eim/internal/config"
)

var (
	Version     string
	ServiceName string
	Branch      string
	Commit      string
	Date        string
)

func Printf() {
	if Version == "" {
		Version = "unknown"
	}
	if ServiceName == "" {
		ServiceName = "EIM-?"
	}
	if Branch == "" {
		Branch = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
	if Date == "" {
		Date = "unknown"
	}

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

	template := `
{{ .Title "` + ServiceName + `" "" 0 }}
{{ .AnsiColor.BrightCyan }}Enterprise Instant Messaging{{ .AnsiColor.Default }}

Version: ` + Version + `
Branch:  ` + Branch + `
Commit:  ` + Commit + `
Date:    ` + strings.Replace(Date, "T", " ", -1) + `

`
	banner.InitString(colorable.NewColorableStdout(), true, true, template)
}
