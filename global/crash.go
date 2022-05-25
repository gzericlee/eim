package global

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"eim/build"
)

func init() {
	if crashFile, err := os.OpenFile(fmt.Sprintf("%v/%v-crash.log", "./logs", strings.ToLower(build.ServiceName)), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664); err == nil {
		crashFile.WriteString(fmt.Sprintf("\n\n\n%v Opened crashfile at %v\n\n", os.Getpid(), time.Now()))
		os.Stderr = crashFile
		syscall.Dup2(int(crashFile.Fd()), 2)
	}
}
