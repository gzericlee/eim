package log

import (
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Field = zap.Field

type Config struct {
	ConsoleEnabled bool
	ConsoleLevel   string
	ConsoleJson    bool

	FileEnabled bool
	FileLevel   string
	FileJson    bool

	Directory  string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func (its *Logger) Debug(msg string, fields ...Field) {
	its.Logger.Debug(msg, fields...)
}

func (its *Logger) Info(msg string, fields ...Field) {
	its.Logger.Info(msg, fields...)
}

func (its *Logger) Warn(msg string, fields ...Field) {
	its.Logger.Warn(msg, fields...)
}

func (its *Logger) Error(msg string, fields ...Field) {
	its.Logger.Error(msg, fields...)
}
func (its *Logger) DPanic(msg string, fields ...Field) {
	its.Logger.DPanic(msg, fields...)
}
func (its *Logger) Panic(msg string, fields ...Field) {
	its.Logger.Panic(msg, fields...)
}
func (its *Logger) Fatal(msg string, fields ...Field) {
	its.Fatal(msg, fields...)
}

type Logger struct {
	*zap.Logger
}

var defaultZapLogger *Logger

func init() {
	defaultZapLogger = newZapLogger(Config{
		ConsoleEnabled: true,
		ConsoleLevel:   "DEBUG",
		ConsoleJson:    false,
		FileEnabled:    false,
		FileJson:       true,
	})
}

func Debug(msg string, fields ...Field) {
	defaultZapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultZapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultZapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	defaultZapLogger.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	defaultZapLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	defaultZapLogger.Fatal(msg, fields...)
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func Configure(config Config) *Logger {
	logger := newZapLogger(config)
	defaultZapLogger = logger
	zap.RedirectStdLog(defaultZapLogger.Logger)
	return logger
}

func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		Error("can't create log directory", zap.Error(err), zap.String("path", config.Directory))
		return nil
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	return zapcore.AddSync(&lumberjack.Logger{
		LocalTime:  true,
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		MaxBackups: config.MaxBackups, // files
	})
}

func newZapLogger(config Config) *Logger {
	var consoleLevel zapcore.Level
	consoleLevel.Set(strings.ToLower(config.ConsoleLevel))

	var fileLevel zapcore.Level
	fileLevel.Set(strings.ToLower(config.FileLevel))

	consoleEncCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "stack_trace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "stack_trace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= consoleLevel
	})
	fileLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= fileLevel
	})

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)
	fileEncoder := zapcore.NewJSONEncoder(jsonEncCfg)

	var cores []zapcore.Core

	if config.ConsoleEnabled {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLevelEnabler))
	}
	if config.FileEnabled {
		cores = append(cores, zapcore.NewCore(fileEncoder, newRollingFile(config), fileLevelEnabler))
	}

	core := zapcore.NewTee(cores...)

	unsugared := zap.New(core).WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{
		Logger: unsugared,
	}
}
