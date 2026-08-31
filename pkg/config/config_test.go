package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("expected sqlite driver, got %q", cfg.Database.Driver)
	}
	want := defaultDatabasePath()
	if cfg.Database.Path != want {
		t.Errorf("expected default path %q, got %q", want, cfg.Database.Path)
	}
}

func TestLoadAppliesDatabaseDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	databasePath := filepath.Join(t.TempDir(), "golem", "cache.db")
	t.Setenv("XDG_CACHE_HOME", filepath.Dir(filepath.Dir(databasePath)))
	err := os.WriteFile(path, []byte(`{"database":{}}`), 0600)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("expected sqlite driver, got %q", cfg.Database.Driver)
	}
	if cfg.Database.Path != databasePath {
		t.Errorf("expected default database path %q, got %q", databasePath, cfg.Database.Path)
	}
}

func TestLoadUsesDefaultConfigPath(t *testing.T) {
	t.Setenv("GOLEM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load default config: %v", err)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("expected sqlite driver, got %q", cfg.Database.Driver)
	}
}

func TestLoadAndSaveSQLiteConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := Config{
		API:      API{Token: "secret", Environment: "sandbox"},
		Database: Database{Driver: "sqlite", Path: defaultDatabasePath()},
	}
	err := Save(path, want)
	if err != nil {
		t.Fatalf("save config: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got != want {
		t.Errorf("loaded config %#v, want %#v", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected config permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestSavePreservesExistingConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := Config{
		API:      API{Token: "old-token", Environment: "sandbox"},
		Database: Database{Driver: "sqlite", Path: "~/custom/cache.db"},
	}
	err := Save(path, want)
	if err != nil {
		t.Fatalf("save initial config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load initial config: %v", err)
	}
	cfg.API.Token = "new-token"
	err = Save(path, cfg)
	if err != nil {
		t.Fatalf("save updated config: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("load updated config: %v", err)
	}
	if got.API.Environment != want.API.Environment {
		t.Errorf("environment = %q, want %q", got.API.Environment, want.API.Environment)
	}
	if got.Database.Path != want.Database.Path {
		t.Errorf("database path = %q, want %q", got.Database.Path, want.Database.Path)
	}
	if got.API.Token != "new-token" {
		t.Errorf("token = %q, want %q", got.API.Token, "new-token")
	}
}

func TestLoadUsesEnvironmentPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("GOLEM_CONFIG", path)
	want := Config{Database: Database{Driver: "sqlite", Path: "cache.db"}}
	err := Save("", want)
	if err != nil {
		t.Fatalf("save config: %v", err)
	}
	got, err := Load("")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if got != want {
		t.Errorf("loaded config %#v, want %#v", got, want)
	}
}

func TestLoadMissingEnvironmentConfigReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	t.Setenv("GOLEM_CONFIG", path)
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load missing environment config: %v", err)
	}
	if cfg != (Default()) {
		t.Errorf("loaded config %#v, want defaults %#v", cfg, Default())
	}
}

func TestLoadExplicitPathOverridesEnvironmentPath(t *testing.T) {
	environmentPath := filepath.Join(t.TempDir(), "environment.json")
	explicitPath := filepath.Join(t.TempDir(), "explicit.json")
	environmentConfig := Config{
		API:      API{Token: "environment-token"},
		Database: Database{Driver: "sqlite", Path: "environment.db"},
	}
	explicitConfig := Config{
		API:      API{Token: "explicit-token"},
		Database: Database{Driver: "sqlite", Path: "explicit.db"},
	}
	err := Save(environmentPath, environmentConfig)
	if err != nil {
		t.Fatalf("save environment config: %v", err)
	}
	err = Save(explicitPath, explicitConfig)
	if err != nil {
		t.Fatalf("save explicit config: %v", err)
	}
	t.Setenv("GOLEM_CONFIG", environmentPath)
	got, err := Load(explicitPath)
	if err != nil {
		t.Fatalf("load explicit config: %v", err)
	}
	if got != explicitConfig {
		t.Errorf("loaded config %#v, want %#v", got, explicitConfig)
	}
}

func TestLoadRejectsUnknownField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := []byte(`{"database":{"driver":"sqlite","path":"cache.db","typo":true}}`)
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{"database":`), 0600)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestLoadRejectsIncompatibleSQLiteFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := []byte(`{"database":{"driver":"sqlite","path":"cache.db","host":"localhost"}}`)
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected incompatible field error")
	}
}

func TestValidateRejectsInvalidAPIEnvironment(t *testing.T) {
	err := validate(Config{API: API{Environment: "invalid"}})
	if err == nil {
		t.Fatal("expected invalid API environment error")
	}
}

func TestValidateRejectsInvalidDatabaseDriver(t *testing.T) {
	err := validate(Config{Database: Database{Driver: "invalid"}})
	if err == nil {
		t.Fatal("expected invalid database driver error")
	}
}

func TestValidateRejectsSQLiteWithoutPath(t *testing.T) {
	err := validate(Config{Database: Database{Driver: "sqlite"}})
	if err == nil {
		t.Fatal("expected missing SQLite path error")
	}
}

func TestValidateRejectsIncompleteServerDatabase(t *testing.T) {
	for _, driver := range []string{"mysql", "postgres"} {
		t.Run(driver, func(t *testing.T) {
			err := validate(Config{Database: Database{Driver: driver}})
			if err == nil {
				t.Fatal("expected missing server database fields error")
			}
		})
	}
}

func TestValidateRejectsServerDatabasePath(t *testing.T) {
	database := Database{
		Driver:   "postgres",
		Path:     "cache.db",
		Host:     "localhost",
		Port:     5432,
		Name:     "database",
		User:     "user",
		Password: "secret",
	}
	err := validate(Config{Database: database})
	if err == nil {
		t.Fatal("expected server database path error")
	}
}

func TestSaveUsesPrivateDirectoryPermissions(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "config", "nested")
	path := filepath.Join(directory, "config.json")
	err := Save(path, Default())
	if err != nil {
		t.Fatalf("save config: %v", err)
	}
	info, err := os.Stat(directory)
	if err != nil {
		t.Fatalf("stat config directory: %v", err)
	}
	if info.Mode().Perm() != 0700 {
		t.Errorf("expected config directory permissions 0700, got %o", info.Mode().Perm())
	}
}

func TestSaveUsesDefaultConfigPath(t *testing.T) {
	t.Setenv("GOLEM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := Save("", Default())
	if err != nil {
		t.Fatalf("save default config: %v", err)
	}
	_, err = os.Stat(defaultConfigPath())
	if err != nil {
		t.Fatalf("stat default config: %v", err)
	}
}

func TestSaveOmitsInapplicableDatabaseFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cfg := Config{
		Database: Database{
			Driver:   "sqlite",
			Path:     "cache.db",
			Host:     "localhost",
			Port:     5432,
			Name:     "database",
			User:     "user",
			Password: "secret",
			SSLMode:  "require",
		},
	}
	err := Save(path, cfg)
	if err != nil {
		t.Fatalf("save sqlite config: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read sqlite config: %v", err)
	}
	fields := []string{"host", "port", "name", "user", "password", "ssl_mode"}
	for _, field := range fields {
		if strings.Contains(string(data), field) {
			t.Errorf("sqlite config contains inapplicable field %q: %s", field, data)
		}
	}

	cfg.Database = Database{Driver: "postgres", Path: "cache.db"}
	err = Save(path, cfg)
	if err != nil {
		t.Fatalf("save postgres config: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read postgres config: %v", err)
	}
	if strings.Contains(string(data), "path") {
		t.Errorf("postgres config contains inapplicable path: %s", data)
	}
}

func TestSaveIncludesConfiguredServerDatabaseFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	database := Database{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		Name:     "database",
		User:     "user",
		Password: "secret",
		SSLMode:  "require",
	}
	cfg := Config{Database: database}
	err := Save(path, cfg)
	if err != nil {
		t.Fatalf("save postgres config: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("load postgres config: %v", err)
	}
	if got != cfg {
		t.Errorf("loaded config %#v, want %#v", got, cfg)
	}
}

func TestSaveAndLoadMySQLConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	database := Database{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Name:     "database",
		User:     "user",
		Password: "secret",
	}
	want := Config{Database: database}
	err := Save(path, want)
	if err != nil {
		t.Fatalf("save mysql config: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("load mysql config: %v", err)
	}
	if got != want {
		t.Errorf("loaded config %#v, want %#v", got, want)
	}
}

func TestExpandPath(t *testing.T) {
	path := ExpandPath("~/cache/golem.db")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get home directory: %v", err)
	}
	want := filepath.Join(home, "cache/golem.db")
	if path != want {
		t.Errorf("expanded path %q, want %q", path, want)
	}
}
