package logger

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func init() {
	logFile = getLogFile()
}

func getLogFile() *os.File {
	currentTime := time.Now()
	logFileName := fmt.Sprintf("logs/log-%d-%02d-%02d.log", currentTime.Year(), currentTime.Month(), currentTime.Day())
	_, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		os.Create(logFileName)
	}
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	return logFile
}

func sendMessage(message string, level string) {
	logFile.WriteString(fmt.Sprintf("%s (%s): %s\n", level, time.Now().Format("2006-01-02 15:04:05"), message))
}

func Info(message string) {
	sendMessage(message, "Info")
}

func Error(message string) {
	sendMessage(message, "Error")
}

func Warning(message string) {
	sendMessage(message, "Warning")
}
