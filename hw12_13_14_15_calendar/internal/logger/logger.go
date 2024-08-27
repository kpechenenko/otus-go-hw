package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var (
	Instance *zap.SugaredLogger
	Level    zap.AtomicLevel
)

const (
	loggerName   = "calendar"
	defaultLevel = zap.InfoLevel
	initLevel    = zap.InfoLevel
)

func init() {
	Level = zap.NewAtomicLevelAt(initLevel)

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
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(os.Stderr), Level)

	Instance = zap.New(core).Sugar().Named(loggerName)
}

func SetLevel(logLevel string) {
	var level zapcore.Level
	switch strings.ToLower(strings.TrimSpace(logLevel)) {
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
	Level.SetLevel(level)
}

func Debug(args ...any) {
	Instance.Debug(args...)
}

func Info(args ...any) {
	Instance.Info(args...)
}

func Warn(args ...any) {
	Instance.Warn(args...)
}

func Error(args ...any) {
	Instance.Error(args...)
}

func Panic(args ...any) {
	Instance.Panic(args...)
}

func Fatal(args ...any) {
	Instance.Fatal(args...)
}

func Debugf(template string, args ...any) {
	Instance.Debugf(template, args...)
}

func Infof(template string, args ...any) {
	Instance.Infof(template, args...)
}

func Warnf(template string, args ...any) {
	Instance.Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	Instance.Errorf(template, args...)
}

func Panicf(template string, args ...any) {
	Instance.Panicf(template, args...)
}

func Fatalf(template string, args ...any) {
	Instance.Fatalf(template, args...)
}

func Errorw(msg string, args ...any) {
	Instance.Errorw(msg, args)
}

func GetNamed(name string) *zap.SugaredLogger {
	return Instance.Named(name)
}
