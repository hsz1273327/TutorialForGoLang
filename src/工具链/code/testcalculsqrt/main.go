package main

import (
	mymath "github.com/tutorialforgolang/calculsqrt/my/mymath"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info(mymath.Sqrt(2))
}
