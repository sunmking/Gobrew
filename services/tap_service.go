package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TapService struct {
	app *application.App
}

func NewTapService(app *application.App) *TapService {
	return &TapService{app: app}
}

func (t *TapService) SetApp(app *application.App) {
	t.app = app
}

func (t *TapService) List(ctx context.Context) ([]TapResult, error) {
	if jsonSupportUnknownOrEnabled(&supportsTapJSONState) {
		stdout, stderr, err := runBrewCommand(ctx, "tap", "--json")
		if err == nil && strings.TrimSpace(stdout) != "" {
			setJSONSupport(&supportsTapJSONState, true)
			var taps []struct {
				Name         string `json:"name"`
				Remote       string `json:"remote"`
				CustomRemote bool   `json:"custom_remote"`
			}
			if err := json.Unmarshal([]byte(stdout), &taps); err != nil {
				return nil, &BrewError{Code: "PARSE_FAILED", Message: "Failed to parse taps", Details: err.Error()}
			}

			results := make([]TapResult, 0, len(taps))
			for _, tap := range taps {
				results = append(results, TapResult{Name: tap.Name, Remote: tap.Remote, CustomRemote: tap.CustomRemote})
			}
			return results, nil
		}
		msg := commandMessage(stdout, stderr)
		if containsUnsupportedJSONOption(msg) {
			setJSONSupport(&supportsTapJSONState, false)
		} else if err != nil && !containsUnsupportedJSONOption(msg) {
			return nil, &BrewError{Code: "TAP_LIST_FAILED", Message: "Failed to list taps", Details: msg}
		}
	}

	textOut, textErr, textExecErr := runBrewCommand(ctx, "tap")
	if textExecErr != nil && strings.TrimSpace(textOut) == "" {
		return nil, &BrewError{Code: "TAP_LIST_FAILED", Message: "Failed to list taps", Details: textErr}
	}

	lines := strings.Split(strings.TrimSpace(textOut), "\n")
	results := make([]TapResult, 0, len(lines))
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		results = append(results, TapResult{Name: name})
	}
	return results, nil
}

func (t *TapService) Add(ctx context.Context, name string) error {
	_, stderr, err := runBrewCommandWithEvents(ctx, t.app, "tap", name)
	if err != nil {
		emitBrewComplete(t.app, false, stderr, 0)
		return &BrewError{Code: "TAP_ADD_FAILED", Message: "Failed to add tap " + name, Details: stderr}
	}
	emitBrewComplete(t.app, true, "", 0)
	return nil
}

func (t *TapService) Remove(ctx context.Context, name string) error {
	_, stderr, err := runBrewCommandWithEvents(ctx, t.app, "untap", name)
	if err != nil {
		emitBrewComplete(t.app, false, stderr, 0)
		return &BrewError{Code: "TAP_REMOVE_FAILED", Message: "Failed to remove tap " + name, Details: stderr}
	}
	emitBrewComplete(t.app, true, "", 0)
	return nil
}

func (t *TapService) Details(ctx context.Context, name string) (*TapDetailResult, error) {
	tapName := strings.TrimSpace(name)
	if tapName == "" {
		return nil, &BrewError{Code: "INVALID_ARGUMENT", Message: "Tap name is required"}
	}

	stdout, stderr, err := runBrewCommand(ctx, "tap-info", "--json", tapName)
	if err != nil {
		return nil, &BrewError{Code: "TAP_INFO_FAILED", Message: "Failed to get tap details", Details: stderr}
	}

	var taps []struct {
		Name         string   `json:"name"`
		Remote       string   `json:"remote"`
		CustomRemote bool     `json:"custom_remote"`
		FormulaNames []string `json:"formula_names"`
		CaskTokens   []string `json:"cask_tokens"`
		LastCommit   string   `json:"last_commit"`
		Branch       string   `json:"branch"`
	}
	if err := json.Unmarshal([]byte(stdout), &taps); err != nil {
		return nil, &BrewError{Code: "PARSE_FAILED", Message: "Failed to parse tap details", Details: err.Error()}
	}

	if len(taps) == 0 {
		return nil, &BrewError{Code: "NOT_FOUND", Message: "Tap not found"}
	}

	selected := taps[0]
	for _, tap := range taps {
		if strings.EqualFold(tap.Name, tapName) {
			selected = tap
			break
		}
	}

	sort.Strings(selected.FormulaNames)
	sort.Strings(selected.CaskTokens)

	return &TapDetailResult{
		Name:         selected.Name,
		Remote:       selected.Remote,
		CustomRemote: selected.CustomRemote,
		FormulaNames: selected.FormulaNames,
		CaskTokens:   selected.CaskTokens,
		LastCommit:   selected.LastCommit,
		Branch:       selected.Branch,
	}, nil
}

func (t *TapService) PackageAction(ctx context.Context, packageType, name, action string) error {
	pkgType := strings.ToLower(strings.TrimSpace(packageType))
	pkgName := strings.TrimSpace(name)
	actionName := strings.ToLower(strings.TrimSpace(action))
	if pkgName == "" {
		return &BrewError{Code: "INVALID_ARGUMENT", Message: "Package name is required"}
	}
	if actionName != "install" && actionName != "uninstall" && actionName != "upgrade" {
		return &BrewError{Code: "INVALID_ARGUMENT", Message: "Unsupported action: " + actionName}
	}
	if pkgType != "formula" && pkgType != "cask" {
		return &BrewError{Code: "INVALID_ARGUMENT", Message: "Unsupported package type: " + pkgType}
	}

	args := []string{actionName}
	if pkgType == "cask" {
		args = append(args, "--cask")
	}
	args = append(args, pkgName)

	_, stderr, err := runBrewCommandWithEvents(ctx, t.app, args...)
	if err != nil {
		emitBrewComplete(t.app, false, stderr, 0)
		return &BrewError{
			Code:    "TAP_PACKAGE_ACTION_FAILED",
			Message: fmt.Sprintf("Failed to %s %s", actionName, pkgName),
			Details: stderr,
		}
	}
	emitBrewComplete(t.app, true, "", 0)
	return nil
}
