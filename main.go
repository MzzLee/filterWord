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
	go signalListen(config)
	ac 	:= libs.AcBuild(config.Keyword)
	server 	:= libs.Create(ac, config)
	server.Start()
}

func signalListen(config *libs.Conf){
	c :=make(chan os.Signal)
	signal.Notify(c)
	select {
	case s := <-c:
		if s == syscall.SIGHUP{
			os.Remove(config.Pid)
			os.Exit(1)
		}
		if s == syscall.SIGQUIT{
			os.Remove(config.Pid)
			os.Exit(1)
		}
		if s == syscall.SIGTERM{
			os.Remove(config.Pid)
			os.Exit(1)
		}

	}
}
