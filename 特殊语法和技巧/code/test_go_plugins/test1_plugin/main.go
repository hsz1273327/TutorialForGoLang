package main

import (
	"fmt"
	"test_go_plugins/shared_module"
)

const plugName string = "test1"

type Test1 struct {
}

func (t *Test1) Sqrt(x float64) float64 {
	shared_module.DoNothing()
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func (t *Test1) Name() string {
	return plugName
}

var Plugin Test1

func init() {
	fmt.Println(plugName, " plugin loading")
	Plugin = Test1{}
}
