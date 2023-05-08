package mylog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	Logger *zap.Logger
)

func InitLogger() {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	Logger, _ = zapConfig.Build(zap.AddCallerSkip(1))
	Logger.Info("Logger init success")
}

func ShutdownLogger() {
	_ = Logger.Sync()
}

func init() {
	InitLogger()
}
