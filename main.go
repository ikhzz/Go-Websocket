package main

import (
	"fmt"
	"clean_arch_v2/config"
	"clean_arch_v2/helper"
	initRouter "clean_arch_v2/helper/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	
	config.InitMysql()

	mw := config.InitFileSetup()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	mid := helper.InitMiddleware()

	router.Use(gin.Recovery())
	router.Use(mid.CORS())
	router.Use(mid.PanicCatcher(mw))
	router.Use(mid.CustomLogger(mw))

	router.LoadHTMLGlob("files/templates/*.html")
	router.Static("/assets", "./files/assets")

	initRouter.InitRouter(router, mw)

	appPort := viper.GetString("address")
	if appPort == "" {
		appPort = ":6060"
	}
	fmt.Println("run at", appPort)
	router.Run(appPort)	
}