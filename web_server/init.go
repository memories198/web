package web_server

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"web/config"
)

var router *gin.Engine

func Start() error {
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	r.Use(gin.Recovery())
	file := config.WebLogFile
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	r.Use(gin.Logger()) //必须写在gin.DefaultWriter后面

	router = r
	router.MaxMultipartMemory = 1 << 30
	registerUrl()
	err := router.Run(":80")
	if err != nil {
		return err
	}
	return nil
}
