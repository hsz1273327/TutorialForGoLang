package config

import (
	"errors"
	logger "testconfigload/logger"

	"github.com/spf13/viper"
)

type config struct {
	Num int
}

//var ConfigViper *viper.Viper

func Init() (config, error) {
	var Config config
	ConfigViper := viper.New()
	InitDefaultConfig(ConfigViper)

	InitFileConfig(ConfigViper)
	InitEnvConfig(ConfigViper)
	InitFlagConfig(ConfigViper)

	err := ConfigViper.Unmarshal(&Config)
	if err != nil {
		logger.Logger.Error("unable to decode into struct, %v", err)
		return Config, errors.New("unable to decode into struct")

	} else {
		if VerifyConfig(Config) {
			return Config, nil
		} else {
			return Config, errors.New("config not satisfied the schema")
		}
	}
}
