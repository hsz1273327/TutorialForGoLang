package config

import (
	"github.com/spf13/viper"
)

func InitEnvConfig(ConfigViper *viper.Viper) {
	EnvConfigViper := viper.New()
	EnvConfigViper.SetEnvPrefix("calcul") // will be uppercased automatically
	EnvConfigViper.BindEnv("num")
	if EnvConfigViper.Get("num") != nil {
		ConfigViper.Set("Num", EnvConfigViper.Get("num"))
	}
}
