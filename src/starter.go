package main

import (
	"clam-server/config"
	"clam-server/serverlogger"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	start()
}

func start() {
	log.Println("clam-server starting ...")
	config.Init()
	serverlogger.Init()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run(":" + config.GetConfig().System.Host)
	if err != nil {
		return
	}
	serverlogger.Warn("clam-server started ...")
}
