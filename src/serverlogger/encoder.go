package serverlogger

import (
	"clam-server/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (c zapcore.EncoderConfig) {
	c = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  config.GetConfig().Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	c.EncodeLevel = zapcore.CapitalColorLevelEncoder

	switch config.GetConfig().Zap.EncodeLevel {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		c.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		c.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		c.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		c.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		c.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return c
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if config.GetConfig().Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getWriteSyncer(fileName) // 日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.GetConfig().Zap.Prefix + config.GetConfig().Zap.TimeFormat))
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,                              // 日志文件的位置
		MaxSize:    config.GetConfig().Zap.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: config.GetConfig().Zap.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     config.GetConfig().Zap.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.GetConfig().Zap.Compress,   // 是否压缩/归档旧文件
	}

	if config.GetConfig().Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
