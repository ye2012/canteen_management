package config

import (
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

var Config = &struct {
	MysqlConfig utils.Config `json:"mysql"`
}{}

func LoadConfig(configPath string) error {
	if err := utils.LoadJSONFile(configPath, Config); err != nil {
		logger.Warn("config", "加载配置错误:%v\n，请检查配置路径: \"%v\"\n,usage: -config 绝对路径\n", err, configPath)
		return err
	}
	logger.Info("config", "Config:%#v", Config)
	return nil
}
