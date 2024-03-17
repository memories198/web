package web_server

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var router *gin.Engine

func Start() error {
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	r.Use(gin.Recovery())
	file, err := os.OpenFile("./logs/web.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	r.Use(gin.Logger()) //必须写在gin.DefaultWriter后面

	router = r
	router.MaxMultipartMemory = 1 << 30
	registerUrl()
	err = router.Run(":80")
	if err != nil {
		return err
	}
	return nil
}
