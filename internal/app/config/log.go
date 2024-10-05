package config

import (
	"log"
	"os"
)

func initLoggers() (infoLogger *log.Logger, errLogger *log.Logger) {
	infoLogger = log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	return
}
