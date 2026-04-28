package services

import (
	"path/filepath"
	"testing"
)

func TestConfigServiceSaveAndLoad(t *testing.T) {
	tempPath := filepath.Join(t.TempDir(), "config.json")
	service := &ConfigService{
		path:   tempPath,
		config: defaultAppConfig(),
	}

	cfg, err := service.Save(AppConfig{
		Language:             "zh",
		Theme:                "dark",
		BrewPath:             "/opt/homebrew/bin/brew",
		BrewfilePath:         "~/Brewfile",
		AutoUpdateInterval:   "6h",
		LogMaxLines:          500,
		CheckUpdatesOnLaunch: false,
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	if cfg.Language != "zh" || cfg.Theme != "dark" || cfg.AutoUpdateInterval != "6h" {
		t.Fatalf("unexpected config after save: %+v", cfg)
	}

	loaded, err := loadConfigFromDisk(tempPath, defaultAppConfig())
	if err != nil {
		t.Fatalf("loadConfigFromDisk returned error: %v", err)
	}
	if loaded.BrewPath != "/opt/homebrew/bin/brew" || loaded.LogMaxLines != 500 {
		t.Fatalf("unexpected loaded config: %+v", loaded)
	}
}

func TestConfigServiceImportSanitizesValues(t *testing.T) {
	tempPath := filepath.Join(t.TempDir(), "config.json")
	service := &ConfigService{
		path:   tempPath,
		config: defaultAppConfig(),
	}

	cfg, err := service.Import(`{"language":"ja","theme":"bad","auto_update_interval":"3h","log_max_lines":5}`)
	if err != nil {
		t.Fatalf("Import returned error: %v", err)
	}

	if cfg.Language != "en" || cfg.Theme != "light" || cfg.AutoUpdateInterval != "off" {
		t.Fatalf("expected sanitized defaults, got: %+v", cfg)
	}
	if cfg.LogMaxLines != 100 {
		t.Fatalf("expected min log max lines, got: %d", cfg.LogMaxLines)
	}
}

func TestConfigServiceReset(t *testing.T) {
	tempPath := filepath.Join(t.TempDir(), "config.json")
	service := &ConfigService{
		path: tempPath,
		config: AppConfig{
			Language: "zh",
			Theme:    "dark",
		},
	}

	cfg, err := service.Reset()
	if err != nil {
		t.Fatalf("Reset returned error: %v", err)
	}
	if cfg.Language != "en" || cfg.Theme != "light" || cfg.LogMaxLines != 1000 {
		t.Fatalf("unexpected reset config: %+v", cfg)
	}
}
