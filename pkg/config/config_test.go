package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Storage.Cache != defaultCachePath() || cfg.Storage.Logs != defaultLogsPath() {
		t.Errorf("storage = %#v, want default storage", cfg.Storage)
	}
}

func TestLoadUsesDefaultConfigPath(t *testing.T) {
	t.Setenv("GOLEM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load default config: %v", err)
	}
	if cfg != Default() {
		t.Fatalf("loaded config %#v, want defaults %#v", cfg, Default())
	}
}

func TestLoadRejectsUnreadableConfigPath(t *testing.T) {
	path := t.TempDir()
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected config read error")
	}
}

func TestLoadAppliesMissingStorageDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{"storage":{}}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Storage.Cache != defaultCachePath() || cfg.Storage.Logs != defaultLogsPath() {
		t.Fatalf("storage = %#v, want defaults", cfg.Storage)
	}
}

func TestLoadRejectsInvalidAPIEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{"api":{"environment":"invalid"}}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected invalid environment error")
	}
}

func TestValidateRequiresStorageFields(t *testing.T) {
	err := validate(Config{Storage: Storage{Logs: "logs"}})
	if err == nil {
		t.Fatal("expected missing cache error")
	}
	err = validate(Config{Storage: Storage{Cache: "cache.db"}})
	if err == nil {
		t.Fatal("expected missing logs error")
	}
}

func TestSaveUsesDefaultConfigPath(t *testing.T) {
	t.Setenv("GOLEM_CONFIG", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	err := Save("", Default())
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(defaultConfigPath())
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadUsesEnvironmentPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("GOLEM_CONFIG", path)
	want := Config{Storage: Storage{Cache: "cache.db", Logs: "logs"}}
	err := Save("", want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("config = %#v, want %#v", got, want)
	}
}

func TestLoadMissingEnvironmentConfigReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	t.Setenv("GOLEM_CONFIG", path)
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load missing environment config: %v", err)
	}
	if cfg != Default() {
		t.Fatalf("loaded config %#v, want defaults %#v", cfg, Default())
	}
}

func TestLoadExplicitPathOverridesEnvironmentPath(t *testing.T) {
	environmentPath := filepath.Join(t.TempDir(), "environment.json")
	explicitPath := filepath.Join(t.TempDir(), "explicit.json")
	environmentStorage := Storage{
		Cache: "environment.db",
		Logs:  "environment-logs",
	}
	explicitStorage := Storage{Cache: "explicit.db", Logs: "explicit-logs"}
	environmentConfig := Config{Storage: environmentStorage}
	explicitConfig := Config{Storage: explicitStorage}
	err := Save(environmentPath, environmentConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = Save(explicitPath, explicitConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOLEM_CONFIG", environmentPath)
	got, err := Load(explicitPath)
	if err != nil {
		t.Fatal(err)
	}
	if got != explicitConfig {
		t.Fatalf("config = %#v, want %#v", got, explicitConfig)
	}
}

func TestSavePreservesExistingConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := Config{
		API:     API{Token: "old-token", Environment: "sandbox"},
		Storage: Storage{Cache: "custom.db", Logs: "logs"},
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
	if got.Storage.Cache != want.Storage.Cache {
		t.Errorf("cache = %q, want %q", got.Storage.Cache, want.Storage.Cache)
	}
	if got.API.Token != "new-token" {
		t.Errorf("token = %q, want %q", got.API.Token, "new-token")
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{"storage":`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestLoadAndSaveConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := Config{
		API:     API{Token: "secret", Environment: "sandbox"},
		Storage: Storage{Cache: "cache.db", Logs: "logs"},
	}
	err := Save(path, want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("config = %#v, want %#v", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestLoadRejectsUnknownField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{"storage":{"cache":"cache.db","unknown":true}}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Load(path)
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestValidateRequiresStoragePaths(t *testing.T) {
	err := validate(Config{})
	if err == nil {
		t.Fatal("expected missing storage error")
	}
}

func TestSaveUsesPrivateDirectoryPermissions(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "config", "nested")
	path := filepath.Join(directory, "config.json")
	err := Save(path, Default())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(directory)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0700 {
		t.Fatalf("permissions = %o, want 700", info.Mode().Perm())
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
