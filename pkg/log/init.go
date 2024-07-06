package log

import (
	"os"
	"syscall"
)

func init() {
	logger = newZapLogger(Config{
		ConsoleEnabled: true,
		ConsoleLevel:   "INFO",
		ConsoleJson:    false,
		FileEnabled:    false,
		FileJson:       true,
	})
}

func SetConfig(cfg Config) {
	logger = newZapLogger(cfg)
}

func New(cfg Config) *Logger {
	return newZapLogger(cfg)
}

func Default() *Logger {
	return logger
}

func RedirectStderr() (err error) {
	logFile, err := os.OpenFile("./crash.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		return
	}
	return
}
