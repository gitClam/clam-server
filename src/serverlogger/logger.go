package serverlogger

import "go.uber.org/zap/zapcore"

func Warn(sg string, fields ...zapcore.Field) {
	z.Warn(sg, fields...)
}

func Info(sg string, fields ...zapcore.Field) {
	z.Info(sg, fields...)
}

func Debug(sg string, fields ...zapcore.Field) {
	z.Debug(sg, fields...)
}

func Error(sg string, fields ...zapcore.Field) {
	z.Error(sg, fields...)
}
