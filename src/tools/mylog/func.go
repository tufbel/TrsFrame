// Package mylog
// Title       : func.go
// Author      : Tuffy  2023/4/3 15:46
// Description :
package mylog

import (
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

func Log(lvl zapcore.Level, msg string, fields ...Field) {
	Logger.Log(lvl, msg, fields...)
}

func Debug(msg string, fields ...Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	Logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	Logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	Logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	Logger.Fatal(msg, fields...)
}
