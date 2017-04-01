package libs

import (
	"log"
	"os"
	"fmt"
	"strings"
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
		createLogPath(logFile)
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

func createLogPath(logFile string) (bool, error){
	logPath := string([]byte(logFile)[0:strings.LastIndex(logFile, "/")])
	status, err := checkPathExists(logPath)
	if status == false || err != nil{
		err := os.MkdirAll(logPath, 0755)
		if err != nil{
			return false, err
		}
		return true, nil
	}
	return false, err
}

func checkPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
