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
	API      API      `mapstructure:"api"`
	Database Database `mapstructure:"database"`
}

type API struct {
	Environment string `mapstructure:"environment"`
	Token       string `mapstructure:"token"`
}

type Database struct {
	Driver   string `mapstructure:"driver"`
	Path     string `mapstructure:"path"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func Default() Config {
	database := Database{Driver: "sqlite", Path: defaultDatabasePath()}
	return Config{Database: database}
}

func defaultDatabasePath() string {
	directory, err := os.UserCacheDir()
	if err != nil {
		//+gocover:ignore:block operating system failure
		return "~/.cache/golem/cache.db"
	}
	return filepath.Join(directory, "golem", "cache.db")
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
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "sqlite"
	}
	if cfg.Database.Driver == "sqlite" && cfg.Database.Path == "" {
		cfg.Database.Path = defaultDatabasePath()
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
	switch strings.ToLower(cfg.Database.Driver) {
	case "sqlite":
		if cfg.Database.Path == "" {
			return fmt.Errorf("database.path is required for sqlite")
		}
		if cfg.Database.Host != "" || cfg.Database.Port != 0 || cfg.Database.Name != "" || cfg.Database.User != "" || cfg.Database.Password != "" || cfg.Database.SSLMode != "" {
			return fmt.Errorf("database host, port, name, user, password and ssl_mode are not valid for sqlite")
		}
	case "mysql", "postgres":
		if cfg.Database.Path != "" {
			return fmt.Errorf("database.path is only valid for sqlite")
		}
		if cfg.Database.Host == "" || cfg.Database.Port == 0 || cfg.Database.Name == "" || cfg.Database.User == "" || cfg.Database.Password == "" {
			return fmt.Errorf("database host, port, name, user and password are required for %s", cfg.Database.Driver)
		}
	default:
		return fmt.Errorf("invalid database driver: %s", cfg.Database.Driver)
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

	databaseSettings := map[string]any{"driver": cfg.Database.Driver}
	switch strings.ToLower(cfg.Database.Driver) {
	case "sqlite":
		databaseSettings["path"] = cfg.Database.Path
	case "mysql", "postgres":
		if cfg.Database.Host != "" {
			databaseSettings["host"] = cfg.Database.Host
		}
		if cfg.Database.Port != 0 {
			databaseSettings["port"] = cfg.Database.Port
		}
		if cfg.Database.Name != "" {
			databaseSettings["name"] = cfg.Database.Name
		}
		if cfg.Database.User != "" {
			databaseSettings["user"] = cfg.Database.User
		}
		if cfg.Database.Password != "" {
			databaseSettings["password"] = cfg.Database.Password
		}
		if cfg.Database.SSLMode != "" {
			databaseSettings["ssl_mode"] = cfg.Database.SSLMode
		}
	}
	v.Set("api", apiSettings)
	v.Set("database", databaseSettings)
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
