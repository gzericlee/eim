package build

import (
	"strings"
	"time"

	"github.com/dimiro1/banner"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
)

var (
	ServiceName string
	Branch      string
	Commit      string
	Date        string
)

func Printf() {
	if ServiceName == "" {
		ServiceName = "EIM-Gateway"
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
	template := `
{{ .Title "` + ServiceName + `" "" 2 }}
   {{ .AnsiColor.BrightCyan }}OpenSource Instant Messaging Server{{ .AnsiColor.Default }}
   Branch: ` + Branch + `
   Commit: ` + Commit + `
   Date:   ` + strings.Replace(Date, "T", " ", -1) + `

`
	banner.InitString(colorable.NewColorableStdout(), true, true, template)
}
