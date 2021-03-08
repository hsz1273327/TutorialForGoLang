package config

import (
	"path"
	"strings"

	"testconfigload/logger"

	"github.com/small-tk/pathlib"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitFlagConfig(ConfigViper *viper.Viper) {
	num := pflag.IntP("num", "n", 0, "要计算平方的值")
	confPath := pflag.StringP("config", "c", "", "配置文件位置")
	pflag.Parse()
	if *confPath != "" {
		p, err := pathlib.New(*confPath).Absolute()
		if err != nil {
			logger.Logger.Info("指定的配置文件获取绝对位置失败")
		} else {
			if p.Exists() && p.IsFile() {
				filenameWithSuffix := path.Base(*confPath)
				fileSuffix := path.Ext(filenameWithSuffix)
				file_name := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
				dir, err := p.Parent()
				if err != nil {
					logger.Logger.Info("指定的配置文件获取父文件夹位置失败")
				} else {
					filePaths := []string{dir.Path}
					SetFileConfig(ConfigViper, file_name, filePaths)
				}

			}
		}
	}
	if *num != 0 {
		ConfigViper.Set("Num", *num) // same result as next line
	}
}
