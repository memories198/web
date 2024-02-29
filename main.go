package main

import (
	"log"
	"web/user"
	web "web/web_server"
)

func main() {
	err := user.DataBaseStart()
	if err != nil {
		log.Println(err)
		return
	}

	err = user.MemoryStart()
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
