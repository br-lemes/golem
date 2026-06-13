package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/viper"
)

var (
	Account    string
	Characters map[string][]string
)

var configDir string

func GetCharacters() []string {
	characters := make([]string, 0, len(Characters))
	for character := range Characters {
		characters = append(characters, character)
	}
	return characters
}

func ConfirmSkill(character string, skill string) bool {
	skills := Characters[character]
	if len(skills) == 0 {
		return true
	}
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	console.Printf("Skill %s is not configured for %s\n", skill, character)
	console.Printf("Available skills: %s\n", strings.Join(skills, ", "))
	return console.Confirm("Do you want to continue?")
}

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configDir = filepath.Join(userConfigDir, "golem")

	token := loadConfig("token")
	api.Token = token.GetString("token")
	if api.Token == "" {
		api.Token = console.Input("Enter your token")
		token.Set("token", api.Token)
		err := token.WriteConfig()
		if err != nil {
			panic(err)
		}
	}

	config := loadConfig("config")
	Account = config.GetString("account")
	if Account == "" {
		details, err := api.MyDetails()
		if err != nil {
			panic(err)
		}
		Account = details.Username
		config.Set("account", Account)
		err = config.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
	Characters = config.GetStringMapStringSlice("characters")
	if len(Characters) != 5 {
		characters, err := api.AccountsCharacters(Account)
		if err != nil {
			panic(err)
		}
		for _, character := range characters {
			Characters[character.Name] = []string{}
		}
		config.Set("characters", Characters)
		err = config.WriteConfig()
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
