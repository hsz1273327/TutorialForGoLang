package logger

import (
	logrus "github.com/sirupsen/logrus"
)

func Init() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

//Logger 默认的logger
var Logger = Init()

//Log 有默认字段的log
var Log = Logger.WithFields(logrus.Fields{
	"common": "this is a common field",
})
