package services

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type ConfigService struct {
	mu     sync.Mutex
	path   string
	config AppConfig
}

func NewConfigService() *ConfigService {
	path := defaultConfigPath()
	config := defaultAppConfig()
	if loaded, err := loadConfigFromDisk(path, config); err == nil {
		config = sanitizeConfig(loaded)
	}
	SetBrewPathOverride(config.BrewPath)
	return &ConfigService{
		path:   path,
		config: config,
	}
}

func (c *ConfigService) Get() AppConfig {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config
}

func (c *ConfigService) Save(input AppConfig) (AppConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	next := sanitizeConfig(input)
	if err := saveConfigToDisk(c.path, next); err != nil {
		return c.config, err
	}
	c.config = next
	SetBrewPathOverride(next.BrewPath)
	return c.config, nil
}

func (c *ConfigService) Reset() (AppConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	next := defaultAppConfig()
	if err := saveConfigToDisk(c.path, next); err != nil {
		return c.config, err
	}
	c.config = next
	SetBrewPathOverride(next.BrewPath)
	return c.config, nil
}

func (c *ConfigService) Export() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	raw, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (c *ConfigService) Import(content string) (AppConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var incoming AppConfig
	if err := json.Unmarshal([]byte(content), &incoming); err != nil {
		return c.config, err
	}
	next := sanitizeConfig(incoming)
	if err := saveConfigToDisk(c.path, next); err != nil {
		return c.config, err
	}
	c.config = next
	SetBrewPathOverride(next.BrewPath)
	return c.config, nil
}

func (c *ConfigService) ValidateBrewPath(path string) BrewPathValidation {
	target := strings.TrimSpace(path)
	if target == "" {
		return BrewPathValidation{
			Path:  "",
			Valid: false,
			Error: "path is empty",
		}
	}

	executable, err := exec.LookPath(target)
	if err != nil {
		return BrewPathValidation{
			Path:  target,
			Valid: false,
			Error: err.Error(),
		}
	}

	out, cmdErr := exec.Command(executable, "--version").Output()
	if cmdErr != nil {
		return BrewPathValidation{
			Path:  executable,
			Valid: false,
			Error: cmdErr.Error(),
		}
	}

	versionLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	return BrewPathValidation{
		Path:    executable,
		Valid:   true,
		Version: versionLine,
	}
}

func (c *ConfigService) ConfigPath() string {
	return c.path
}

func defaultAppConfig() AppConfig {
	return AppConfig{
		Language:             "en",
		Theme:                "light",
		BrewPath:             "",
		BrewfilePath:         "",
		AutoUpdateInterval:   "off",
		LogMaxLines:          1000,
		CheckUpdatesOnLaunch: true,
	}
}

func defaultConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, "Gobrew", "config.json")
	}
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "gobrew-config.json"
	}
	return filepath.Join(homeDir, ".gobrew", "config.json")
}

func sanitizeConfig(input AppConfig) AppConfig {
	cfg := input
	if cfg.Language != "zh" && cfg.Language != "en" {
		cfg.Language = "en"
	}
	if cfg.Theme != "dark" && cfg.Theme != "light" {
		cfg.Theme = "light"
	}
	switch cfg.AutoUpdateInterval {
	case "off", "1h", "6h", "12h", "24h":
	default:
		cfg.AutoUpdateInterval = "off"
	}
	if cfg.LogMaxLines < 100 {
		cfg.LogMaxLines = 100
	}
	if cfg.LogMaxLines > 10000 {
		cfg.LogMaxLines = 10000
	}
	cfg.BrewPath = strings.TrimSpace(cfg.BrewPath)
	cfg.BrewfilePath = strings.TrimSpace(cfg.BrewfilePath)
	return cfg
}

func loadConfigFromDisk(path string, fallback AppConfig) (AppConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fallback, err
	}
	next := fallback
	if err := json.Unmarshal(raw, &next); err != nil {
		return fallback, err
	}
	return next, nil
}

func saveConfigToDisk(path string, cfg AppConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0644)
}
