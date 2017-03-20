package libs

import (
	"log"
	"os"
	"fmt"
)

var (
	ENV_DEBUG  = "DEBUG"
	ENV_TEST   = "TEST"
	ENV_ONLINE = "ONLINE"
)

var (
	_loggerInstance *GLogger
)

type GLogger struct {
	Logger *log.Logger
	Env string
}

func InitLogger(logFile string, Env string) (*GLogger , error){
	if _loggerInstance == nil {
		handler, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		_, err = handler.Seek(0, os.SEEK_END)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		_loggerInstance = new(GLogger)
		_loggerInstance.Logger 	= log.New(handler, "", log.Ldate | log.Ltime | log.Llongfile)
		_loggerInstance.Env	= Env
	}
	return _loggerInstance, nil
}

func (gl *GLogger) Fatal(v ...interface{}) {
	gl.Logger.SetPrefix("FATAL ")
	gl.Logger.Println(v)
}

func (gl *GLogger) Warning(v ...interface{}) {
	if gl.Env == ENV_DEBUG || gl.Env == ENV_TEST {
		gl.Logger.SetPrefix("WARNING ")
		gl.Logger.Println(v)
	}
}

func (gl *GLogger) Notice(v ...interface{}) {
	if gl.Env == ENV_DEBUG || gl.Env == ENV_TEST {
		gl.Logger.SetPrefix("NOTICE ")
		gl.Logger.Println(v)
	}
}

func (gl *GLogger) Info(v ...interface{}) {
	if gl.Env == ENV_DEBUG {
		gl.Logger.SetPrefix("INFO ")
		gl.Logger.Println(v)
	}
}

func (gl *GLogger) Trace(v ...interface{}) {
	if gl.Env == ENV_DEBUG {
		gl.Logger.SetPrefix("TRACE ")
		gl.Logger.Println(v)
	}
}
