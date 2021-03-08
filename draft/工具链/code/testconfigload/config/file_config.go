package config

import (
	"testconfigload/logger"

	"github.com/spf13/viper"
)

func SetFileConfig(ConfigViper *viper.Viper, file_name string, filePaths []string) {
	FileConfigViper := viper.New()
	FileConfigViper.SetConfigName(file_name)
	for _, file_path := range filePaths {
		FileConfigViper.AddConfigPath(file_path)
	}
	err := FileConfigViper.ReadInConfig() // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		logger.Logger.Info("config file not found: %s \n", err)
	} else {
		logger.Logger.Info("Num in file is %d \n", FileConfigViper.Get("Num"))
		ConfigViper.Set("Num", FileConfigViper.Get("Num"))
		logger.Logger.Info(ConfigViper.Get("Num"))
	}
}

func InitFileConfig(ConfigViper *viper.Viper) {
	file_name := "config"
	filePaths := []string{"/etc/appname/", "$HOME/.appname", "."}
	SetFileConfig(ConfigViper, file_name, filePaths)
}
