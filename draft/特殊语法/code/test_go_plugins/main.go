package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
	"test_go_plugins/shared_module"
)

type PluginInterface interface {
	Sqrt(float64) float64
	Name() string
}

var Plugins map[string]PluginInterface = map[string]PluginInterface{}

func main() {
	shared_module.DoNothing()
	// 加载动态库
	rd, err := ioutil.ReadDir("plugins")
	if err != nil {
		fmt.Printf("read dir get error: %v \n", err)
		os.Exit(1)
	}
	for _, fi := range rd {
		file_name := fi.Name()
		fmt.Println(fi.Name())
		if fi.IsDir() {
			fmt.Printf("[%s] is dir skip\n", file_name)
		} else {
			if strings.HasSuffix(file_name, "plugin") {
				module_path := fmt.Sprintf("plugins/%s", file_name)
				module, err := plugin.Open(module_path)
				if err != nil {
					fmt.Println("plugin ", file_name, "load error: ", err.Error())

				} else {
					plugin_s, err := module.Lookup("Plugin")
					if err != nil {
						fmt.Println("plugin ", file_name, "Lookup Plugin error: ", err.Error())
					} else {
						plugin := plugin_s.(PluginInterface)
						Plugins[plugin.Name()] = plugin
						fmt.Println("load plugin ", plugin.Name(), " from ", module_path, " ok")
					}
				}
			}
		}
	}
	for plugin_name, plugin := range Plugins {
		fmt.Println("plugin ", plugin_name, "run sqrt(2) get result", plugin.Sqrt(2))
	}
}

func init() {
	fmt.Println("running main")
}
