package config

import (
	"github.com/spf13/viper"
)

func InitDefaultConfig(ConfigViper *viper.Viper) {
	ConfigViper.SetDefault("Num", 2)
}
