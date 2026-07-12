package api

import (
	"os"
	"path/filepath"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/viper"
)

var configDir string

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configDir = filepath.Join(userConfigDir, "golem")

	tokenConfig := loadConfig("token")
	token = tokenConfig.GetString("token")
	if token == "" {
		token = console.Input("Enter your token")
		tokenConfig.Set("token", token)
		err := tokenConfig.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func loadConfig(name string) *viper.Viper {
	configFile := filepath.Join(configDir, name+".json")
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0755)
		err := viper.WriteConfigAs(configFile)
		if err != nil {
			panic(err)
		}
	}
	config := viper.New()
	config.AddConfigPath(configDir)
	config.SetConfigName(name)
	config.SetConfigType("json")
	err = config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return config
}
