package config

import (
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

var Config = &struct {
	MysqlConfig      utils.Config `json:"mysql"`
	FileStorePath    string       `json:"file_store_path"`
	MchID            string       `json:"mch_id"`
	MchCertSerialNum string       `json:"mch_cert_serial_num"`
	MchApiV3Key      string       `json:"mch_api_v3_key"`
	PrivateKeyPath   string       `json:"private_key_path"`
	AppID            string       `json:"app_id"`
	AppSecret        string       `json:"app_secret"`
	KitchenAppID     string       `json:"kitchen_app_id"`
	KitchenAppSecret string       `json:"kitchen_app_secret"`
}{}

func LoadConfig(configPath string) error {
	if err := utils.LoadJSONFile(configPath, Config); err != nil {
		logger.Warn("config", "加载配置错误:%v\n，请检查配置路径: \"%v\"\n,usage: -config 绝对路径\n", err, configPath)
		return err
	}
	logger.Info("config", "Config:%#v", Config)
	return nil
}

const CustomKey = "Custom"
const TokenKey = "Token"
