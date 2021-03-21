package main

import (
	"QuickPass/crontab"
	"QuickPass/models"
	"QuickPass/pkg/gredis"
	"QuickPass/pkg/logf"
	"QuickPass/pkg/mq"
	"QuickPass/pkg/setting"
	"QuickPass/pkg/util"
	vali "QuickPass/pkg/validation"
	"QuickPass/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
	logf.Setup()
	//upload.SetUp()
	models.Setup()
	util.Setup()
	gredis.Setup()
	mq.SocketChanSetup()
	crontab.Setup()
	vali.InitValidation()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://Pay/blob/master/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
