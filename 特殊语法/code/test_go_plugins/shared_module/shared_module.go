package shared_module

import "fmt"

func DoNothing() {
	fmt.Println("do nothing")
}

func init() {
	fmt.Println("shared module init1")
}
func init() {
	fmt.Println("shared module init2")
}
func init() {
	fmt.Println("shared module init3")
}
