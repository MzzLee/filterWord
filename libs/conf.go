package libs

import (
	"flag"
	"github.com/larspensjo/config"
	"log"
	"fmt"
	"sync"
)

type Conf struct {
	File string
	Bind string
	Port int
	Protocol string
	Keyword string
	Pid string
	LogFile string
}

var (
	_confInstance *Conf
	_lock *sync.Mutex = &sync.Mutex{}
)

func (conf *Conf) Argv() *Conf {
	conf.File      = *flag.String("c","conf/config.ini", "General configuration file")
	conf.Keyword   = *flag.String("key","", "Sensitive word file")
	conf.LogFile   = *flag.String("log","", "Log file")
	flag.Parse()
	return conf
}

func (conf *Conf) Load() *Conf {

	conf.Argv()

	buffer, err := config.ReadDefault(conf.File)
	if err != nil{
		log.Fatalf("Fail to Find : %s %s", conf.File, err)
	}

	conf.Bind, _ 		= buffer.String("Server", "Bind")
	conf.Port, _ 		= buffer.Int("Server", "Port")
	conf.Protocol, _ 	= buffer.String("Server", "Protocol")
	conf.Pid,_              = buffer.String("Server", "Pid")
	if conf.Keyword == ""{
		conf.Keyword, _ = buffer.String("Server", "Keyword")
	}
	if conf.LogFile == ""{
		conf.LogFile, _ = buffer.String("Server", "LogFile")
	}

	return conf
}

func GetConfigInstance() *Conf{
	if _confInstance == nil{
		_lock.Lock()
		defer _lock.Unlock()

		_confInstance = new(Conf)
		_confInstance.Load()
	}
	return _confInstance
}

