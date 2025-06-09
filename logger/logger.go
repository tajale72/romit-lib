package logger

import (
	"log"
	"os"
)

type RomitLogger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	traceLogger *log.Logger
}

type RomitLoggerInterface interface {
	LogInfo(message string)
	LogWarning(message string)
	LogError(message string)
	LogFatal(message string)
	LogPanic(message string)
	LogDebug(message string)
	LogTrace(message string)
}

func NewRomitLogger() RomitLoggerInterface {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	return &RomitLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", flags),
		warnLogger:  log.New(os.Stdout, "WARN: ", flags),
		errorLogger: log.New(os.Stdout, "ERROR: ", flags),
	}
}

func (l *RomitLogger) LogInfo(message string) {
	l.infoLogger.Output(2, message)
}

func (l *RomitLogger) LogWarning(message string) {
	l.warnLogger.Output(2, message)
}

func (l *RomitLogger) LogError(message string) {
	l.errorLogger.Output(2, message)
}

func (l *RomitLogger) LogFatal(message string) {
	l.errorLogger.Fatal(message)
}

func (l *RomitLogger) LogPanic(message string) {
	l.errorLogger.Output(2, message)
}

func (l *RomitLogger) LogDebug(message string) {
	l.debugLogger.Output(2, message)
}

func (l *RomitLogger) LogTrace(message string) {
	l.traceLogger.Output(2, message)
}
