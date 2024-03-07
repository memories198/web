package main

import (
	"log"
	"web/config"
	"web/dao"
	web "web/web_server"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Println(err)
		return
	}
	err = dao.DataBaseStart()
	if err != nil {
		log.Println(err)
		return
	}

	err = dao.MemoryStart()
	if err != nil {
		log.Println(err)
		return
	}

	err = web.Start()
	if err != nil {
		log.Println(err)
		return
	}

}
