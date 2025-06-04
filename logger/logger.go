package logger

import (
	"log"
	"os"
)

type RomitLogger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
}

func NewRomitLogger() *RomitLogger {
	return &RomitLogger{
		infoLogger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger: log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *RomitLogger) LogInfo(message string) {
	l.infoLogger.Println(message)
}

func (l *RomitLogger) LogWarning(message string) {
	l.warnLogger.Println(message)
}
