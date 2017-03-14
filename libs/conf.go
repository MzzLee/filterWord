package libs

import (
	"flag"
	"github.com/larspensjo/config"
	"log"
)

type Conf struct {
	Bind string
	Port int
	Protocol string
	Keyword string
}

var _confInstance *Conf

func (c *Conf)Load() *Conf {

	configFile := flag.String("c","conf/config.ini", "General configuration file")
	keyFile    := flag.String("key","", "General configuration file")
	flag.Parse()

	buffer, err := config.ReadDefault(*configFile)
	if err != nil{
		log.Fatalf("Fail to Find : %s %s", *configFile, err)
	}
	c.Keyword 		= *keyFile
	c.Bind, _ 		= buffer.String("Server", "Bind")
	c.Port, _ 		= buffer.Int("Server", "Port")
	c.Protocol, _ 		= buffer.String("Server", "Protocol")
	if c.Keyword == ""{
		c.Keyword, _ 		= buffer.String("Server", "Keyword")
	}

	return c
}

func ConfInstance() *Conf{
	if _confInstance == nil{
		_confInstance = new(Conf).Load()
	}
	return _confInstance
}

