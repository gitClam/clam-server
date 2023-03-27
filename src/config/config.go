package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	path     = "../config/config.yaml"
	fileType = "yaml"
)

var v *viper.Viper
var c *config

func GetConfig() *config {
	return c
}

func Init() {
	log.Println("init config starting ...")
	v = viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(fileType)
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("load config file err : %s\n", err)
		os.Exit(2)
	}
	v.OnConfigChange(onConfigChange)
	v.WatchConfig()
	if err := v.Unmarshal(&c); err != nil {
		log.Printf("load config err : %s\n", err)
		os.Exit(2)
	}
	log.Println("init config started ...")
}

func onConfigChange(e fsnotify.Event) {
	log.Println("config reload:", e.Name)
	var newConfig *config
	if err := v.Unmarshal(newConfig); err != nil {
		log.Printf("config reload err: %s\n", e)
		return
	}
	c = newConfig
}
