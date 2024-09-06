package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var instance *zap.SugaredLogger

const (
	defaultLevel = zap.InfoLevel
	loggerName   = "calendar"
)

func init() {
	initGlobalLoggerWithLevel(defaultLevel)
}

func initGlobalLoggerWithLevel(level zapcore.Level) {
	config := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			zapcore.RFC3339NanoTimeEncoder(time.UTC(), encoder)
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(os.Stderr), level)
	instance = zap.New(core).
		Sugar().
		Named(loggerName)

	instance.Infof("Logger initialized with level: %s", level)
}

func SetLevel(l string) {
	var level zapcore.Level
	switch strings.ToLower(l) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = defaultLevel
	}
	initGlobalLoggerWithLevel(level)
}

func Debug(args ...any) {
	instance.Debug(args...)
}

func Info(args ...any) {
	instance.Info(args...)
}

func Warn(args ...any) {
	instance.Warn(args...)
}

func Error(args ...any) {
	instance.Error(args...)
}

func Panic(args ...any) {
	instance.Panic(args...)
}

func Fatal(args ...any) {
	instance.Fatal(args...)
}

func Debugf(template string, args ...any) {
	instance.Debugf(template, args...)
}

func Infof(template string, args ...any) {
	instance.Infof(template, args...)
}

func Warnf(template string, args ...any) {
	instance.Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	instance.Errorf(template, args...)
}

func Panicf(template string, args ...any) {
	instance.Panicf(template, args...)
}

func Fatalf(template string, args ...any) {
	instance.Fatalf(template, args...)
}
