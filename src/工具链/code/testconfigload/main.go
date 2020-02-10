package main

import (
	"testconfigload/config"
	"testconfigload/logger"

	mymath "github.com/tutorialforgolang/calculsqrt/my/mymath"
)

func main() {
	conf, err := config.Init()
	if err != nil {
		logger.Logger.Info(err)
	} else {
		logger.Logger.Info(conf)
		logger.Logger.Info(mymath.Sqrt(float64(conf.Num)))
	}
}
