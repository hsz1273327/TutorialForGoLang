package config

import (
	"os"

	"github.com/paked/configure"
)

// Config 配置项
type Config struct {
	Address string
	Debug   bool
}

// LoadConfig 从文件,环境变量,参数中加载配置项
func LoadConfig() Config {
	conf := configure.New()
	address := conf.String("addr", "localhost:8080", "http service address")
	debug := conf.Bool("debug", false, "debug mod or not")
	_, err := os.Stat("config.json")
	if err != nil {
		if os.IsExist(err) {
			conf.Use(configure.NewJSONFromFile("config.json"))
		}
	}
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewFlagWithUsage(usage))
	conf.Parse()
	return Config{Address: *address, Debug: *debug}

}

func usage() string {
	return `简单的用法说明
	--addr=localhost:8080  http service address
	--debug=false          debug mod or not
	`
}
