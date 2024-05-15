package log

func init() {
	logger = newZapLogger(Config{
		ConsoleEnabled: true,
		ConsoleLevel:   "DEBUG",
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
