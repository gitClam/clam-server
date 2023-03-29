package main

import (
	"clam-server/config"
	"clam-server/cors"
	"clam-server/jwt"
	"clam-server/serverlogger"
	"clam-server/service/logdecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	start()
}

func start() {
	log.Println("clam-server starting ...")
	initServerBase()
	initServerComponents()
	r := gin.New()
	initGinComponents(r)
	initRouter(r)
	serverHeart(r)
	err := r.Run(":" + config.GetConfig().System.Host)
	if err != nil {
		serverlogger.Warn("start clam-server fail", zap.String("err", err.Error()))
		return
	}
	serverlogger.Warn("clam-server started ...")
}

func initServerBase() {
	config.Init()
	serverlogger.Init()
}

func initServerComponents() {

}

func initGinComponents(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(serverlogger.LoggerHandler())
	r.Use(cors.Cors())
	r.Use(jwts.JwtHandler())
}

func initRouter(r *gin.Engine) {
	logdecode.Router(r)
}

func serverHeart(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
