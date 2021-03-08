package logger

import (
	logrus "github.com/sirupsen/logrus"
)

//Init 初始化logger
func Init() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

//Logger 默认的logger
var Logger = Init()

var contextLogger = Logger.WithFields(logrus.Fields{
	"name": "ws_test"})

// Trace 默认logger的Trace级别的log
func Trace(args ...interface{}) {
	contextLogger.Trace(args...)
}

// Debug 默认logger的Debug级别的log
func Debug(args ...interface{}) {
	contextLogger.Debug(args...)
}

// Info 默认logger的Info级别的log
func Info(args ...interface{}) {
	contextLogger.Info(args...)
}

// Warn 默认logger的Warn级别的log
func Warn(args ...interface{}) {
	contextLogger.Warn(args...)
}

// Error 默认logger的Error级别的log
func Error(args ...interface{}) {
	contextLogger.Error(args...)
}

// Fatal 默认logger的Fatal级别的log
func Fatal(args ...interface{}) {
	contextLogger.Fatal(args...)
}

// Panic 默认logger的Panic级别的log
func Panic(args ...interface{}) {
	contextLogger.Panic(args...)
}
