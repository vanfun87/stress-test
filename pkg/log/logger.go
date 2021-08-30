package log

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	logFilepath  = "./stress-test.log"
	EnableLogger bool
)

func InitLogger() {
	if !EnableLogger {
		return
	}

	log.SetFormatter(&log.JSONFormatter{})

	var file, err = os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
	}

	log.SetOutput(file)

	log.SetLevel(log.InfoLevel)
}

func Info(args ...interface{}) {
	if !EnableLogger {
		return
	}

	log.Info(args...)
}

func Println(args ...interface{}) {
	if !EnableLogger {
		return
	}

	log.Println(args...)
}

func Printf(str string, args ...interface{}) {
	if !EnableLogger {
		return
	}

	log.Printf(str, args...)
}
