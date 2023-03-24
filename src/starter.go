package main

import (
	"clam-server/config"
	"clam-server/serverlogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func main() {
	start()
}

func start() {
	log.Println("clam-server starting ...")
	config.Init()
	serverlogger.Init()
	serverlogger.Warn("start clam-server fail", zap.String("a", "abc"))
	r := gin.New()
	gin.Default()
	r.Use(gin.Recovery())
	r.Use(serverlogger.GinLogger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run(":" + config.GetConfig().System.Host)
	if err != nil {
		serverlogger.Warn("start clam-server fail", zap.String("err", err.Error()))
		return
	}
	serverlogger.Warn("clam-server started ...")
}
