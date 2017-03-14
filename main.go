package main

import "./libs"

func main(){
	config 	:= libs.ConfInstance()
	ac 	:= libs.AcBuild(config.Keyword)
	server 	:= libs.Create(ac, config)
	server.Start()
}
