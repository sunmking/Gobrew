package services

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type BundleService struct {
	app *application.App
}

func NewBundleService(app *application.App) *BundleService {
	return &BundleService{app: app}
}

func (b *BundleService) SetApp(app *application.App) {
	b.app = app
}

func (b *BundleService) Dump(ctx context.Context, filePath string, force bool) error {
	resolved, err := resolveBrewfilePath(filePath)
	if err != nil {
		return err
	}

	args := []string{"bundle", "dump", "--file", resolved, "--describe"}
	if force {
		args = append(args, "--force")
	}

	start := time.Now()
	_, stderr, runErr := runBrewCommandWithEvents(ctx, b.app, args...)
	b.emitComplete(runErr == nil, stderr, start)
	if runErr != nil {
		return &BrewError{Code: "BUNDLE_DUMP_FAILED", Message: "Failed to dump Brewfile", Details: stderr}
	}
	return nil
}

func (b *BundleService) Restore(ctx context.Context, filePath string) error {
	args := []string{"bundle", "install", "--verbose"}
	if strings.TrimSpace(filePath) != "" {
		args = append(args, "--file", filePath)
	}

	start := time.Now()
	_, stderr, runErr := runBrewCommandWithEvents(ctx, b.app, args...)
	b.emitComplete(runErr == nil, stderr, start)
	if runErr != nil {
		return &BrewError{Code: "BUNDLE_RESTORE_FAILED", Message: "Failed to restore from Brewfile", Details: stderr}
	}
	return nil
}

func (b *BundleService) Check(ctx context.Context, filePath string) (*BundleCheckResult, error) {
	args := []string{"bundle", "check", "--verbose"}
	if strings.TrimSpace(filePath) != "" {
		args = append(args, "--file", filePath)
	}

	stdout, stderr, err := runBrewCommand(ctx, args...)
	if err != nil {
		missing := parseMissingDependencies(stdout, stderr)
		return &BundleCheckResult{Satisfied: false, Missing: missing, Output: stdout + stderr}, nil
	}
	return &BundleCheckResult{Satisfied: true, Missing: nil, Output: stdout}, nil
}

func (b *BundleService) CleanupPreview(ctx context.Context, filePath string) (*CommandResult, error) {
	// `brew bundle cleanup` has no --dry-run flag.
	// Running without --force acts as a preview and exits non-zero if removals are needed.
	args := []string{"bundle", "cleanup", "--verbose"}
	if strings.TrimSpace(filePath) != "" {
		args = append(args, "--file", filePath)
	}

	start := time.Now()
	stdout, stderr, err := runBrewCommand(ctx, args...)
	output := strings.TrimSpace(stdout + "\n" + stderr)
	if err != nil && output == "" {
		return nil, &BrewError{Code: "BUNDLE_CLEANUP_PREVIEW_FAILED", Message: "Failed to preview bundle cleanup", Details: stderr}
	}
	return &CommandResult{
		Success:  err == nil,
		Output:   output,
		Duration: time.Since(start).String(),
		Error:    stderr,
	}, nil
}

func (b *BundleService) Cleanup(ctx context.Context, filePath string, force bool) error {
	args := []string{"bundle", "cleanup", "--verbose"}
	if strings.TrimSpace(filePath) != "" {
		args = append(args, "--file", filePath)
	}
	if force {
		args = append(args, "--force")
	}

	start := time.Now()
	_, stderr, runErr := runBrewCommandWithEvents(ctx, b.app, args...)
	b.emitComplete(runErr == nil, stderr, start)
	if runErr != nil {
		return &BrewError{Code: "BUNDLE_CLEANUP_FAILED", Message: "Failed to cleanup bundle", Details: stderr}
	}
	return nil
}

func (b *BundleService) List(ctx context.Context, filePath string) ([]string, error) {
	args := []string{"bundle", "list", "--all"}
	if strings.TrimSpace(filePath) != "" {
		args = append(args, "--file", filePath)
	}

	stdout, stderr, err := runBrewCommand(ctx, args...)
	if err != nil && stdout == "" {
		return nil, &BrewError{Code: "BUNDLE_LIST_FAILED", Message: "Failed to list bundle", Details: stderr}
	}

	result := make([]string, 0)
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func (b *BundleService) ReadBrewfile(filePath string) (string, error) {
	resolved, err := resolveBrewfilePath(filePath)
	if err != nil {
		return "", err
	}

	data, readErr := os.ReadFile(resolved)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return "", &BrewError{Code: "FILE_NOT_FOUND", Message: "Brewfile not found", Details: resolved}
		}
		return "", &BrewError{Code: "READ_FAILED", Message: "Failed to read Brewfile", Details: readErr.Error()}
	}
	return string(data), nil
}

func (b *BundleService) emitComplete(success bool, details string, start time.Time) {
	emitBrewComplete(b.app, success, details, time.Since(start))
}

func parseMissingDependencies(stdout, stderr string) []string {
	missing := make([]string, 0)
	for _, line := range strings.Split(stdout+stderr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "is not installed") || strings.Contains(strings.ToLower(line), "missing") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				missing = append(missing, fields[0])
			}
		}
	}
	return missing
}

func resolveBrewfilePath(filePath string) (string, error) {
	if strings.TrimSpace(filePath) != "" {
		return filePath, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", &BrewError{Code: "PATH_FAILED", Message: "Failed to get home directory", Details: err.Error()}
	}
	return filepath.Join(homeDir, "Brewfile"), nil
}
