package libs

import (
	"log"
	"os"
	"fmt"
)

var (
	_loggerInstance *log.Logger
)

func InitLogger(logFile string) (* log.Logger, error){
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
		_loggerInstance = log.New(handler, "", log.Ldate | log.Ltime | log.Llongfile)
	}
	return _loggerInstance, nil
}
