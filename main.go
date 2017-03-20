package main

import (
	"./libs"
)

func main(){

	status := libs.DeamonClient().Run()
	if status  == false {
		return
	}
	server 	:= libs.CreateServer()
	server.Start()
}

