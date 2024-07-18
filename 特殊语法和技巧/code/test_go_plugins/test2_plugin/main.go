package main

import (
	"fmt"
	"test_go_plugins/shared_module"
)

const plugName string = "test2"

type Test2 struct {
}

func (t *Test2) Sqrt(x float64) float64 {
	shared_module.DoNothing()
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
func (t *Test2) Name() string {
	return plugName
}

var Plugin Test2

func init() {
	fmt.Println(plugName, " plugin loading")
	Plugin = Test2{}
}
