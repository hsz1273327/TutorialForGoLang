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
