package tieba

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

var configDirPath string
var configFilePath string

func InitConfig() error {
	userPath, _ := os.UserHomeDir()
	configDirPath = path.Join(userPath, ".tieEpub")
	err := os.MkdirAll(configDirPath, 0755)
	if err != nil {
		return err
	}
	configFilePath = path.Join(configDirPath, "config.json")
	viper.SetConfigName("config")
	viper.AddConfigPath(configDirPath)
	viper.SetDefault("cookie", "")
	viper.SetDefault("ai-api-key", "")
	_ = viper.ReadInConfig()
	err = viper.WriteConfigAs(configFilePath)
	if err != nil {
		return err
	}
	return nil
}
func SaveConfig() error {
	return viper.WriteConfigAs(configFilePath)
}
