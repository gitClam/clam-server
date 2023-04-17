package main

import (
	"clam-server/component/cors"
	"clam-server/component/jwt"
	"clam-server/config"
	"clam-server/serverlogger"
	"clam-server/service/logdecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strconv"
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
	writePid()
	err := r.Run(":" + config.GetConfig().System.Host)
	if err != nil {
		serverlogger.Warn("start clam-server fail", zap.String("err", err.Error()))
		return
	}
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

func writePid() {
	f, err := os.Create("../pid")
	if err != nil {
		return
	}
	_, err = f.Write([]byte(strconv.Itoa(os.Getpid())))
	if err != nil {
		return
	}
}
