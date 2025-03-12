package utils

import (
	"log"
	"os"
)

// InitLogger initializes the logging configuration.
func InitLogger() {
	logFile, err := os.OpenFile("exporter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// LogInfo logs informational messages.
func LogInfo(message string) {
	log.Println("INFO: " + message)
}