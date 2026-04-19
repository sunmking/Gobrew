package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type BrewService struct {
	app *application.App
}

var (
	supportsListJSONV2State   int32 = -1 // -1 unknown, 0 unsupported, 1 supported
	supportsOutdatedJSONState int32 = -1
	supportsTapJSONState      int32 = -1
)

func NewBrewService(app *application.App) *BrewService {
	return &BrewService{app: app}
}

func (b *BrewService) SetApp(app *application.App) {
	b.app = app
}

func (b *BrewService) Capabilities(ctx context.Context) BrewCapabilities {
	return getBrewCapabilities(ctx)
}

func (b *BrewService) ListInstalled(ctx context.Context) (*InstalledListResult, error) {
	if jsonSupportUnknownOrEnabled(&supportsListJSONV2State) {
		formulaJSON, formulaErrText, formulaErr := runBrewCommand(ctx, "list", "--formula", "--json=v2")
		caskJSON, caskErrText, caskErr := runBrewCommand(ctx, "list", "--cask", "--json=v2")
		if formulaErr == nil && caskErr == nil {
			setJSONSupport(&supportsListJSONV2State, true)
			formulae, err := parseFormulaListJSON(formulaJSON)
			if err != nil {
				return nil, &BrewError{Code: "PARSE_FAILED", Message: "Failed to parse formulae list", Details: err.Error()}
			}
			casks, err := parseCaskListJSON(caskJSON)
			if err != nil {
				return nil, &BrewError{Code: "PARSE_CASK_FAILED", Message: "Failed to parse cask list", Details: err.Error()}
			}
			return &InstalledListResult{Formulae: formulae, Casks: casks}, nil
		}
		formulaMsg := commandMessage(formulaJSON, formulaErrText)
		caskMsg := commandMessage(caskJSON, caskErrText)
		if containsUnsupportedJSONOption(formulaMsg) || containsUnsupportedJSONOption(caskMsg) {
			setJSONSupport(&supportsListJSONV2State, false)
		} else if formulaErr != nil && !containsUnsupportedJSONOption(formulaMsg) {
			return nil, &BrewError{Code: "LIST_FAILED", Message: "Failed to list formulae", Details: formulaMsg}
		} else if caskErr != nil && !containsUnsupportedJSONOption(caskMsg) {
			return nil, &BrewError{Code: "LIST_CASK_FAILED", Message: "Failed to list casks", Details: caskMsg}
		}
	}

	// Fast compatibility path for older Homebrew: plain list + versions.
	type cmdResult struct {
		stdout string
		stderr string
		err    error
	}
	chFormula := make(chan cmdResult, 1)
	chCask := make(chan cmdResult, 1)
	chFormulaVer := make(chan cmdResult, 1)
	chCaskVer := make(chan cmdResult, 1)

	go func() {
		stdout, stderr, err := runBrewCommand(ctx, "list", "--formula")
		chFormula <- cmdResult{stdout: stdout, stderr: stderr, err: err}
	}()
	go func() {
		stdout, stderr, err := runBrewCommand(ctx, "list", "--cask")
		chCask <- cmdResult{stdout: stdout, stderr: stderr, err: err}
	}()
	go func() {
		stdout, stderr, err := runBrewCommand(ctx, "list", "--versions", "--formula")
		chFormulaVer <- cmdResult{stdout: stdout, stderr: stderr, err: err}
	}()
	go func() {
		stdout, stderr, err := runBrewCommand(ctx, "list", "--versions", "--cask")
		chCaskVer <- cmdResult{stdout: stdout, stderr: stderr, err: err}
	}()

	formulaRes := <-chFormula
	caskRes := <-chCask
	formulaVerRes := <-chFormulaVer
	caskVerRes := <-chCaskVer

	formulaText, formulaErrText, formulaErr := formulaRes.stdout, formulaRes.stderr, formulaRes.err
	caskText, caskErrText, caskErr := caskRes.stdout, caskRes.stderr, caskRes.err
	if formulaErr != nil && strings.TrimSpace(formulaText) == "" {
		return nil, &BrewError{Code: "LIST_FAILED", Message: "Failed to list formulae", Details: formulaErrText}
	}
	if caskErr != nil && strings.TrimSpace(caskText) == "" {
		return nil, &BrewError{Code: "LIST_CASK_FAILED", Message: "Failed to list casks", Details: caskErrText}
	}

	formulaVersionsText := formulaVerRes.stdout
	caskVersionsText := caskVerRes.stdout
	formulaVersions := parseVersionsMap(formulaVersionsText)
	caskVersions := parseVersionsMap(caskVersionsText)

	return &InstalledListResult{
		Formulae: parseFormulaListText(formulaText, formulaVersions),
		Casks:    parseCaskListText(caskText, caskVersions),
	}, nil
}

func (b *BrewService) Search(ctx context.Context, query string) (*SearchResult, error) {
	formulaeOut, formulaeErr, formulaeExecErr := runBrewCommand(ctx, "search", "--formula", query)
	casksOut, casksErr, casksExecErr := runBrewCommand(ctx, "search", "--cask", query)

	if formulaeExecErr != nil && strings.TrimSpace(formulaeOut) == "" && casksExecErr != nil && strings.TrimSpace(casksOut) == "" {
		return nil, &BrewError{
			Code:    "SEARCH_FAILED",
			Message: "Search failed",
			Details: strings.TrimSpace(formulaeErr + "\n" + casksErr),
		}
	}

	toItems := func(raw string) []SearchItem {
		items := make([]SearchItem, 0)
		for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			items = append(items, SearchItem{Name: line, FullName: line})
		}
		return items
	}

	return &SearchResult{
		Formulae: toItems(formulaeOut),
		Casks:    toItems(casksOut),
	}, nil
}

func (b *BrewService) Info(ctx context.Context, name string) (*FormulaInfo, error) {
	stdout, stderr, err := runBrewCommand(ctx, "info", "--json=v2", name)
	if err != nil {
		return nil, &BrewError{Code: "INFO_FAILED", Message: "Failed to get info for " + name, Details: stderr}
	}

	var raw struct {
		Formulae []struct {
			Name      string    `json:"name"`
			FullName  string    `json:"full_name"`
			Tap       string    `json:"tap"`
			Desc      string    `json:"desc"`
			Homepage  string    `json:"homepage"`
			License   string    `json:"license"`
			LinkedKeg string    `json:"linked_keg"`
			Pinned    bool      `json:"pinned"`
			Installed []KegInfo `json:"installed"`
			Versions  struct {
				Stable string `json:"stable"`
				Head   string `json:"head"`
			} `json:"versions"`
			Dependencies []string `json:"dependencies"`
		} `json:"formulae"`
	}

	if err := json.Unmarshal([]byte(stdout), &raw); err != nil {
		return nil, &BrewError{Code: "PARSE_FAILED", Message: "Failed to parse info", Details: err.Error()}
	}

	if len(raw.Formulae) == 0 {
		return nil, &BrewError{Code: "NOT_FOUND", Message: "Package not found: " + name}
	}

	f := raw.Formulae[0]
	return &FormulaInfo{
		Name:         f.Name,
		FullName:     f.FullName,
		Tap:          f.Tap,
		Desc:         f.Desc,
		Homepage:     f.Homepage,
		License:      f.License,
		StableVer:    f.Versions.Stable,
		HeadVer:      f.Versions.Head,
		LinkedKeg:    f.LinkedKeg,
		Pinned:       f.Pinned,
		Installed:    f.Installed,
		Dependencies: f.Dependencies,
	}, nil
}

func (b *BrewService) Install(ctx context.Context, name string) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "install", name)
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "INSTALL_FAILED", Message: "Failed to install " + name, Details: stderr}
	}
	return nil
}

func (b *BrewService) Uninstall(ctx context.Context, name string) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "uninstall", name)
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "UNINSTALL_FAILED", Message: "Failed to uninstall " + name, Details: stderr}
	}
	return nil
}

func (b *BrewService) UninstallCask(ctx context.Context, name string, zap bool) error {
	start := time.Now()
	args := []string{"uninstall", "--cask", name}
	if zap {
		args = append(args, "--zap")
	}
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, args...)
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "UNINSTALL_CASK_FAILED", Message: "Failed to uninstall cask " + name, Details: stderr}
	}
	return nil
}

func (b *BrewService) Upgrade(ctx context.Context, name string) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "upgrade", name)
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "UPGRADE_FAILED", Message: "Failed to upgrade " + name, Details: stderr}
	}
	return nil
}

func (b *BrewService) UpgradeAll(ctx context.Context) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "upgrade")
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "UPGRADE_ALL_FAILED", Message: "Failed to upgrade all packages", Details: stderr}
	}
	return nil
}

func (b *BrewService) Update(ctx context.Context) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "update")
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "UPDATE_FAILED", Message: "Failed to update Homebrew", Details: stderr}
	}
	return nil
}

func (b *BrewService) Reinstall(ctx context.Context, name string, isCask bool) error {
	start := time.Now()
	args := []string{"reinstall"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, args...)
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "REINSTALL_FAILED", Message: "Failed to reinstall " + name, Details: stderr}
	}
	return nil
}

func (b *BrewService) Deps(ctx context.Context, name string, tree bool) ([]string, error) {
	args := []string{"deps"}
	if tree {
		args = append(args, "--tree")
	}
	args = append(args, name)
	stdout, stderr, err := runBrewCommand(ctx, args...)
	if err != nil && strings.TrimSpace(stdout) == "" {
		return nil, &BrewError{Code: "DEPS_FAILED", Message: "Failed to list dependencies for " + name, Details: stderr}
	}
	return parseNameLines(stdout), nil
}

func (b *BrewService) Uses(ctx context.Context, name string, installedOnly, recursive bool) ([]string, error) {
	args := []string{"uses"}
	if installedOnly {
		args = append(args, "--installed")
	}
	if recursive {
		args = append(args, "--recursive")
	}
	args = append(args, name)
	stdout, stderr, err := runBrewCommand(ctx, args...)
	if err != nil && strings.TrimSpace(stdout) == "" {
		return nil, &BrewError{Code: "USES_FAILED", Message: "Failed to list dependents for " + name, Details: stderr}
	}
	return parseNameLines(stdout), nil
}

func (b *BrewService) Outdated(ctx context.Context) (*OutdatedListResult, error) {
	if jsonSupportUnknownOrEnabled(&supportsOutdatedJSONState) {
		stdout, stderr, err := runBrewCommand(ctx, "outdated", "--json=v2")
		if err == nil && strings.TrimSpace(stdout) != "" {
			setJSONSupport(&supportsOutdatedJSONState, true)
			var raw struct {
				Formulae []OutdatedFormula `json:"formulae"`
				Casks    []OutdatedCask    `json:"casks"`
			}
			if unmarshalErr := json.Unmarshal([]byte(stdout), &raw); unmarshalErr != nil {
				return nil, &BrewError{Code: "PARSE_FAILED", Message: "Failed to parse outdated list", Details: unmarshalErr.Error()}
			}
			return &OutdatedListResult{Formulae: raw.Formulae, Casks: raw.Casks}, nil
		}
		if strings.TrimSpace(stdout) == "" && err == nil {
			return &OutdatedListResult{}, nil
		}
		msg := commandMessage(stdout, stderr)
		if containsUnsupportedJSONOption(msg) {
			setJSONSupport(&supportsOutdatedJSONState, false)
		} else if err != nil && !containsUnsupportedJSONOption(msg) {
			return nil, &BrewError{Code: "OUTDATED_FAILED", Message: "Failed to list outdated packages", Details: msg}
		}
	}

	formulaText, _, _ := runBrewCommand(ctx, "outdated", "--formula")
	caskText, _, _ := runBrewCommand(ctx, "outdated", "--cask")
	return &OutdatedListResult{
		Formulae: parseOutdatedFormulaText(formulaText),
		Casks:    parseOutdatedCaskText(caskText),
	}, nil
}

func jsonSupportUnknownOrEnabled(state *int32) bool {
	v := atomic.LoadInt32(state)
	return v == -1 || v == 1
}

func setJSONSupport(state *int32, supported bool) {
	if supported {
		atomic.StoreInt32(state, 1)
		return
	}
	atomic.StoreInt32(state, 0)
}

func (b *BrewService) CleanupPreview(ctx context.Context) (*CleanupResult, error) {
	stdout, _, _ := runBrewCommand(ctx, "cleanup", "--dry-run")
	lineCount := countNonEmptyLines(stdout)
	return &CleanupResult{CleanedCount: lineCount, Output: stdout}, nil
}

func (b *BrewService) Cleanup(ctx context.Context) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "cleanup")
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "CLEANUP_FAILED", Message: "Cleanup failed", Details: stderr}
	}
	return nil
}

func (b *BrewService) AutoRemovePreview(ctx context.Context) (*CleanupResult, error) {
	stdout, stderr, err := runBrewCommand(ctx, "autoremove", "--dry-run")
	if err != nil && strings.TrimSpace(stdout) == "" {
		return nil, &BrewError{Code: "AUTOREMOVE_PREVIEW_FAILED", Message: "Autoremove preview failed", Details: stderr}
	}
	return &CleanupResult{CleanedCount: countNonEmptyLines(stdout), Output: stdout}, nil
}

func (b *BrewService) AutoRemove(ctx context.Context) error {
	start := time.Now()
	_, stderr, err := runBrewCommandWithEvents(ctx, b.app, "autoremove")
	b.emitComplete(err == nil, stderr, start)
	if err != nil {
		return &BrewError{Code: "AUTOREMOVE_FAILED", Message: "Autoremove failed", Details: stderr}
	}
	return nil
}

func (b *BrewService) Doctor(ctx context.Context) (*CommandResult, error) {
	start := time.Now()
	stdout, stderr, err := runBrewCommand(ctx, "doctor")
	return &CommandResult{
		Success:  err == nil,
		Output:   stdout,
		Error:    stderr,
		Duration: time.Since(start).String(),
	}, nil
}

func (b *BrewService) emitComplete(success bool, details string, start time.Time) {
	if b.app == nil {
		return
	}
	duration := time.Since(start).String()
	if success {
		b.app.Event.Emit("brew-complete", fmt.Sprintf(`{"success":true,"duration":"%s"}`, duration))
		return
	}
	esc := strings.ReplaceAll(details, `"`, `'`)
	b.app.Event.Emit("brew-complete", fmt.Sprintf(`{"success":false,"error":"%s","duration":"%s"}`, esc, duration))
}

func parseFormulaListJSON(stdout string) ([]FormulaInfo, error) {
	var rawFormula struct {
		Formulae []struct {
			Name      string `json:"name"`
			FullName  string `json:"full_name"`
			Tap       string `json:"tap"`
			Desc      string `json:"desc"`
			Homepage  string `json:"homepage"`
			License   string `json:"license"`
			LinkedKeg string `json:"linked_keg"`
			Pinned    bool   `json:"pinned"`
			Installed []struct {
				Version               string `json:"version"`
				InstalledAsDependency bool   `json:"installed_as_dependency"`
				InstalledOnRequest    bool   `json:"installed_on_request"`
			} `json:"installed"`
			Versions struct {
				Stable string `json:"stable"`
				Head   string `json:"head"`
			} `json:"versions"`
			Dependencies []string `json:"dependencies"`
		} `json:"formulae"`
	}
	if err := json.Unmarshal([]byte(stdout), &rawFormula); err != nil {
		return nil, err
	}
	formulae := make([]FormulaInfo, 0, len(rawFormula.Formulae))
	for _, f := range rawFormula.Formulae {
		installed := make([]KegInfo, 0, len(f.Installed))
		for _, inst := range f.Installed {
			installed = append(installed, KegInfo{
				Version:               inst.Version,
				InstalledAsDependency: inst.InstalledAsDependency,
				InstalledOnRequest:    inst.InstalledOnRequest,
			})
		}
		formulae = append(formulae, FormulaInfo{
			Name:         f.Name,
			FullName:     f.FullName,
			Tap:          f.Tap,
			Desc:         f.Desc,
			Homepage:     f.Homepage,
			License:      f.License,
			StableVer:    f.Versions.Stable,
			HeadVer:      f.Versions.Head,
			LinkedKeg:    f.LinkedKeg,
			Pinned:       f.Pinned,
			Installed:    installed,
			Dependencies: f.Dependencies,
		})
	}
	return formulae, nil
}

func parseCaskListJSON(stdout string) ([]CaskInfo, error) {
	var rawCask struct {
		Casks []struct {
			Name        string `json:"name"`
			FullName    string `json:"full_name"`
			Tap         string `json:"tap"`
			Desc        string `json:"desc"`
			Homepage    string `json:"homepage"`
			Version     string `json:"version"`
			Installed   string `json:"installed"`
			AutoUpdates bool   `json:"auto_updates"`
			Token       string `json:"token"`
		} `json:"casks"`
	}
	if err := json.Unmarshal([]byte(stdout), &rawCask); err != nil {
		return nil, err
	}
	casks := make([]CaskInfo, 0, len(rawCask.Casks))
	for _, c := range rawCask.Casks {
		casks = append(casks, CaskInfo{
			Name:        c.Name,
			FullName:    c.FullName,
			Tap:         c.Tap,
			Desc:        c.Desc,
			Homepage:    c.Homepage,
			Version:     c.Version,
			Installed:   c.Installed,
			AutoUpdates: c.AutoUpdates,
			Token:       c.Token,
		})
	}
	return casks, nil
}

func parseInstalledFromInfoJSON(stdout string) ([]FormulaInfo, []CaskInfo, error) {
	var raw struct {
		Formulae []map[string]any `json:"formulae"`
		Casks    []map[string]any `json:"casks"`
	}
	if err := json.Unmarshal([]byte(stdout), &raw); err != nil {
		return nil, nil, err
	}

	formulae := make([]FormulaInfo, 0, len(raw.Formulae))
	for _, f := range raw.Formulae {
		stable, head := "", ""
		if versions, ok := f["versions"].(map[string]any); ok {
			stable = stringValue(versions["stable"])
			head = stringValue(versions["head"])
		}
		name := stringValue(f["name"])
		fullName := stringValue(f["full_name"])
		if fullName == "" {
			fullName = name
		}
		formulae = append(formulae, FormulaInfo{
			Name:         name,
			FullName:     fullName,
			Tap:          stringValue(f["tap"]),
			Desc:         stringValue(f["desc"]),
			Homepage:     stringValue(f["homepage"]),
			License:      stringValue(f["license"]),
			StableVer:    stable,
			HeadVer:      head,
			LinkedKeg:    stringValue(f["linked_keg"]),
			Pinned:       boolValue(f["pinned"]),
			Installed:    parseKegInfos(f["installed"]),
			Dependencies: stringSliceValue(f["dependencies"]),
		})
	}

	casks := make([]CaskInfo, 0, len(raw.Casks))
	for _, c := range raw.Casks {
		name := stringValue(c["name"])
		token := stringValue(c["token"])
		if token == "" {
			token = name
		}
		fullName := stringValue(c["full_name"])
		if fullName == "" {
			fullName = token
		}
		casks = append(casks, CaskInfo{
			Name:        name,
			FullName:    fullName,
			Tap:         stringValue(c["tap"]),
			Desc:        stringValue(c["desc"]),
			Homepage:    stringValue(c["homepage"]),
			Version:     stringValue(c["version"]),
			Installed:   installedValue(c["installed"]),
			AutoUpdates: boolValue(c["auto_updates"]),
			Token:       token,
		})
	}
	return formulae, casks, nil
}

func parseFormulaListText(stdout string, versions map[string]string) []FormulaInfo {
	items := make([]FormulaInfo, 0)
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		version := versions[name]
		installed := make([]KegInfo, 0, 1)
		if version != "" {
			installed = append(installed, KegInfo{Version: version})
		}
		items = append(items, FormulaInfo{
			Name:      name,
			FullName:  name,
			StableVer: version,
			Installed: installed,
		})
	}
	return items
}

func parseCaskListText(stdout string, versions map[string]string) []CaskInfo {
	items := make([]CaskInfo, 0)
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		items = append(items, CaskInfo{
			Name:      name,
			FullName:  name,
			Token:     name,
			Version:   versions[name],
			Installed: versions[name],
		})
	}
	return items
}

func parseVersionsMap(stdout string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		result[fields[0]] = fields[1]
	}
	return result
}

func parseOutdatedFormulaText(stdout string) []OutdatedFormula {
	items := make([]OutdatedFormula, 0)
	for _, name := range parseNameLines(stdout) {
		items = append(items, OutdatedFormula{Name: name})
	}
	return items
}

func parseOutdatedCaskText(stdout string) []OutdatedCask {
	items := make([]OutdatedCask, 0)
	for _, name := range parseNameLines(stdout) {
		items = append(items, OutdatedCask{Name: name})
	}
	return items
}

func parseNameLines(stdout string) []string {
	items := make([]string, 0)
	for _, line := range strings.Split(stdout, "\n") {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		items = append(items, name)
	}
	return items
}

func countNonEmptyLines(stdout string) int {
	return len(parseNameLines(stdout))
}

func parseKegInfos(v any) []KegInfo {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	items := make([]KegInfo, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		items = append(items, KegInfo{
			Version:               stringValue(m["version"]),
			InstalledAsDependency: boolValue(m["installed_as_dependency"]),
			InstalledOnRequest:    boolValue(m["installed_on_request"]),
		})
	}
	return items
}

func containsUnsupportedJSONOption(text string) bool {
	msg := strings.ToLower(text)
	return strings.Contains(msg, "invalid option: --json=v2") || strings.Contains(msg, "invalid option: --json")
}

func stringValue(v any) string {
	switch x := v.(type) {
	case string:
		return x
	default:
		return ""
	}
}

func boolValue(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	default:
		return false
	}
}

func stringSliceValue(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	items := make([]string, 0, len(raw))
	for _, item := range raw {
		s := stringValue(item)
		if s == "" {
			continue
		}
		items = append(items, s)
	}
	return items
}

func installedValue(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []any:
		items := make([]string, 0, len(x))
		for _, item := range x {
			s := stringValue(item)
			if s != "" {
				items = append(items, s)
			}
		}
		return strings.Join(items, ", ")
	default:
		return ""
	}
}
