package logger

import (
	log "github.com/sirupsen/logrus"
)

var singletonLogger log.Logger

func Init(level log.Level) {
	singletonLogger = *log.New()
	singletonLogger.SetLevel(level)
}

func Debug(args ...interface{}) {
	singletonLogger.Debug(args...)
}

func Info(args ...interface{}) {
	singletonLogger.Info(args...)
}

func Error(args ...interface{}) {
	singletonLogger.Error(args...)
}
