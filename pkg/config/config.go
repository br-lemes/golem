package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	API     API     `mapstructure:"api"`
	Storage Storage `mapstructure:"storage"`
}

type API struct {
	Environment string `mapstructure:"environment"`
	Token       string `mapstructure:"token"`
}

type Storage struct {
	Cache string `mapstructure:"cache"`
	Logs  string `mapstructure:"logs"`
}

func Default() Config {
	storage := Storage{Cache: defaultCachePath(), Logs: defaultLogsPath()}
	return Config{Storage: storage}
}

func defaultCachePath() string {
	directory, err := os.UserCacheDir()
	if err != nil {
		//+gocover:ignore:block operating system failure
		return "~/.cache/golem/cache.db"
	}
	return filepath.Join(directory, "golem", "cache.db")
}

func defaultLogsPath() string {
	directory, err := os.UserCacheDir()
	if err != nil {
		//+gocover:ignore:block operating system failure
		return "~/.cache/golem/logs"
	}
	return filepath.Join(directory, "golem", "logs")
}

func defaultConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		//+gocover:ignore:block operating system failure
		return "~/.config/golem/config.json"
	}
	return filepath.Join(dir, "golem", "config.json")
}

func Load(path string) (Config, error) {
	if path == "" {
		path = os.Getenv("GOLEM_CONFIG")
		if path == "" {
			path = defaultConfigPath()
		}
	}
	path = ExpandPath(path)

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	err = v.UnmarshalExact(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}
	if cfg.Storage.Cache == "" {
		cfg.Storage.Cache = defaultCachePath()
	}
	if cfg.Storage.Logs == "" {
		cfg.Storage.Logs = defaultLogsPath()
	}
	err = validate(cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validate(cfg Config) error {
	switch strings.ToLower(cfg.API.Environment) {
	case "", "standard", "sandbox", "beta":
	default:
		return fmt.Errorf("invalid api environment: %s", cfg.API.Environment)
	}
	if cfg.Storage.Cache == "" {
		return fmt.Errorf("storage.cache is required")
	}
	if cfg.Storage.Logs == "" {
		return fmt.Errorf("storage.logs is required")
	}
	return nil
}

func Save(path string, cfg Config) error {
	if path == "" {
		path = os.Getenv("GOLEM_CONFIG")
		if path == "" {
			path = defaultConfigPath()
		}
	}
	path = ExpandPath(path)
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		//+gocover:ignore:block operating system failure
		return fmt.Errorf("create config directory: %w", err)
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	apiSettings := map[string]any{}
	if cfg.API.Token != "" {
		apiSettings["token"] = cfg.API.Token
	}
	if cfg.API.Environment != "" {
		apiSettings["environment"] = cfg.API.Environment
	}

	storageSettings := map[string]any{
		"cache": cfg.Storage.Cache,
		"logs":  cfg.Storage.Logs,
	}
	v.Set("api", apiSettings)
	v.Set("storage", storageSettings)
	err = v.WriteConfigAs(path)
	if err != nil {
		//+gocover:ignore:block operating system failure
		return fmt.Errorf("write config file: %w", err)
	}
	return os.Chmod(path, 0600)
}

func ExpandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return os.ExpandEnv(path)
}
