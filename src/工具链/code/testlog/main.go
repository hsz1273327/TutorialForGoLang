package main

import (
	"testlog/logger"

	logrus "github.com/sirupsen/logrus"
)

func main() {
	logger.Logger.SetLevel(logrus.InfoLevel)
	logger.Logger.Info("测试")
	logger.Log.Info("测试")
	logger.Logger.WithFields(logrus.Fields{
		"event": "field",
	}).Info("测试 field")
	logger.Log.WithFields(logrus.Fields{
		"event": "field",
	}).Info("测试 field")

	logger.Logger.SetLevel(logrus.WarnLevel)
	logger.Logger.Info("测试 INFO")
	logger.Log.Info("测试 INFO")
	logger.Logger.Warn("测试 warn")
	logger.Log.Warn("测试 warn")
}
