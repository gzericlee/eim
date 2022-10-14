package log

func init() {
	defaultZapLogger = newZapLogger(Config{
		ConsoleEnabled: true,
		ConsoleLevel:   "DEBUG",
		ConsoleJson:    false,
		FileEnabled:    false,
		FileJson:       true,
	})
}

func InitLogger(cfg Config) {
	defaultZapLogger = newZapLogger(cfg)
}
