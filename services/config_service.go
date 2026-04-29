package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type ConfigService struct {
	mu             sync.Mutex
	path           string
	config         AppConfig
	menuBarApplier func(bool)
}

func NewConfigService() *ConfigService {
	path := defaultConfigPath()
	config := defaultAppConfig()
	if loaded, err := loadConfigFromDisk(path, config); err == nil {
		config = sanitizeConfig(loaded)
	}
	SetBrewPathOverride(config.BrewPath)
	SetCommandRuntimeOptions(config.MaxConcurrency, config.DebugLog)
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
	SetCommandRuntimeOptions(next.MaxConcurrency, next.DebugLog)
	c.applyRuntimeEffectsLocked()
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
	SetCommandRuntimeOptions(next.MaxConcurrency, next.DebugLog)
	c.applyRuntimeEffectsLocked()
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
	SetCommandRuntimeOptions(next.MaxConcurrency, next.DebugLog)
	c.applyRuntimeEffectsLocked()
	return c.config, nil
}

func (c *ConfigService) ValidateBrewPath(path string) BrewPathValidation {
	target := strings.TrimSpace(path)
	if target == "" {
		return BrewPathValidation{
			Path:    "",
			Valid:   true,
			Version: "auto-detect",
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

func (c *ConfigService) SetMenuBarApplier(applier func(bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.menuBarApplier = applier
	c.applyMenuBarLocked(c.config.ShowInMenuBar)
}

func defaultAppConfig() AppConfig {
	return AppConfig{
		Language:                "en",
		Theme:                   "light",
		BrewPath:                "",
		BrewfilePath:            "",
		AutoUpdateInterval:      "off",
		LogMaxLines:             1000,
		CheckUpdatesOnLaunch:    true,
		LaunchAtLogin:           false,
		ShowInMenuBar:           true,
		RestoreLastPage:         true,
		LastPagePath:            "/",
		BackupBeforeUpdate:      true,
		CleanupAfterUpdate:      false,
		NotifyUpdates:           true,
		NotifyOperations:        true,
		NotifyServices:          true,
		NotifyErrors:            true,
		ShowToasts:              true,
		ToastPosition:           "br",
		KeepNotificationHistory: true,
		ToastDurationMs:         2500,
		MaxConcurrency:          4,
		DebugLog:                false,
		AccentHue:               255,
		UIFontSize:              13,
		RowDensity:              "default",
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
	if cfg.ToastDurationMs < 1000 {
		cfg.ToastDurationMs = 1000
	}
	if cfg.ToastDurationMs > 10000 {
		cfg.ToastDurationMs = 10000
	}
	if cfg.ToastPosition == "" {
		cfg.ToastPosition = "br"
	}
	switch cfg.ToastPosition {
	case "tl", "tr", "bl", "br":
	default:
		cfg.ToastPosition = "br"
	}
	if cfg.MaxConcurrency < 1 {
		cfg.MaxConcurrency = 1
	}
	if cfg.MaxConcurrency > 16 {
		cfg.MaxConcurrency = 16
	}
	if cfg.LastPagePath == "" {
		cfg.LastPagePath = "/"
	}
	if cfg.UIFontSize < 12 {
		cfg.UIFontSize = 12
	}
	if cfg.UIFontSize > 16 {
		cfg.UIFontSize = 16
	}
	if cfg.AccentHue < 0 || cfg.AccentHue > 360 {
		cfg.AccentHue = 255
	}
	switch cfg.RowDensity {
	case "compact", "default", "comfortable":
	default:
		cfg.RowDensity = "default"
	}
	cfg.BrewPath = strings.TrimSpace(cfg.BrewPath)
	cfg.BrewfilePath = strings.TrimSpace(cfg.BrewfilePath)
	cfg.LastPagePath = strings.TrimSpace(cfg.LastPagePath)
	return cfg
}

func (c *ConfigService) EffectiveBrewPath() (string, error) {
	path, err := resolveBrewPath()
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *ConfigService) ValidateStartupConfig() (string, error) {
	cfg := c.Get()
	return fmt.Sprintf("launch_at_login=%t show_in_menu_bar=%t", cfg.LaunchAtLogin, cfg.ShowInMenuBar), nil
}

func (c *ConfigService) applyRuntimeEffectsLocked() {
	_ = applyLaunchAtLogin(c.config.LaunchAtLogin)
	c.applyMenuBarLocked(c.config.ShowInMenuBar)
}

func (c *ConfigService) applyMenuBarLocked(enabled bool) {
	if IsDevRuntime() {
		return
	}
	if c.menuBarApplier == nil {
		return
	}
	c.menuBarApplier(enabled)
}

func applyLaunchAtLogin(enabled bool) error {
	if IsDevRuntime() {
		return nil
	}
	if runtime.GOOS != "darwin" {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	agentDir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return err
	}
	agentPath := filepath.Join(agentDir, "com.gobrew.app.plist")
	if !enabled {
		_ = runLaunchCtl("bootout", "gui/"+currentUID(), agentPath)
		if _, statErr := os.Stat(agentPath); statErr == nil {
			return os.Remove(agentPath)
		}
		return nil
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.gobrew.app</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <false/>
</dict>
</plist>`, executable)
	if err := os.WriteFile(agentPath, []byte(content), 0644); err != nil {
		return err
	}
	_ = runLaunchCtl("bootstrap", "gui/"+currentUID(), agentPath)
	_ = runLaunchCtl("enable", "gui/"+currentUID()+"/com.gobrew.app")
	return nil
}

func currentUID() string {
	return fmt.Sprintf("%d", os.Getuid())
}

func runLaunchCtl(args ...string) error {
	cmd := exec.Command("launchctl", args...)
	return cmd.Run()
}

func IsDevRuntime() bool {
	// Wails dev / Vite dev server markers.
	envMarkers := []string{
		"WAILS_DEV_SERVER",
		"WAILS_DEV_SERVER_URL",
		"WAILS_VITE_DEV_SERVER",
		"VITE_DEV_SERVER_URL",
	}
	for _, key := range envMarkers {
		if strings.TrimSpace(os.Getenv(key)) != "" {
			return true
		}
	}
	// Common command markers when running under dev tooling.
	args := os.Args
	return slices.Contains(args, "dev") || slices.Contains(args, "wails3") || slices.Contains(args, "wails")
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
