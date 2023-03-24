package serverlogger

import (
	"clam-server/config"
	"clam-server/utils/io"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var z *zap.Logger

func Init() {
	log.Println("init serverLogger starting ...")
	zapDirector := config.GetConfig().Zap.Director
	// 判断是否有文件夹
	if !io.CheckDirExists(zapDirector) {
		log.Printf("create directory %s\n", zapDirector)
		if err := os.Mkdir(zapDirector, os.ModePerm); err != nil {
			log.Printf("create directory fail %s\n", err)
			os.Exit(2)
		}
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", zapDirector), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", zapDirector), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", zapDirector), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", zapDirector), errorPriority),
	}
	z = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	if config.GetConfig().Zap.ShowLine {
		z = z.WithOptions(zap.AddCaller())
	}
	log.Println("init serverLogger started ...")
}
