package web_server

import (
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"web/config"
)

var router *gin.Engine
var serverKey []byte

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

	k, err := os.ReadFile("certificate/server.key")
	if err != nil {
		return err
	}
	block, rest := pem.Decode(k)
	if len(rest) != 0 {
		return err
	}

	serverKey = block.Bytes

	err = router.Run(":80")
	if err != nil {
		return err
	}
	return nil
}
