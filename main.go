package main

import (
	"os"
	"os/signal"
	"syscall"

	"./libs"
)

func main(){
	config 	:= libs.ConfInstance()
	daemon  := libs.Daemon{Config:config}
	status  := daemon.Run()

	if status  == false {
		return
	}
	go signalListen()
	ac 	:= libs.AcBuild(config.Keyword)
	server 	:= libs.Create(ac, config)
	server.Start()
}

func signalListen(){
	c :=make(chan os.Signal)
	signal.Notify(c)
	select {
	case s := <-c:
		if s == syscall.SIGHUP{
			os.Exit(1)
		}
		if s == syscall.SIGQUIT{
			os.Exit(1)
		}
		if s == syscall.SIGTERM{
			os.Exit(1)
		}

	}
}
