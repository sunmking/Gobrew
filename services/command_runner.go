package services

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	brewPathOnce sync.Once
	brewPath     string
	brewPathErr  error
	brewPathMu   sync.RWMutex
	brewPathUser string
	runtimeMu    sync.Mutex
	brewSem      chan struct{}
	debugLog     bool
)

func resolveBrewPath() (string, error) {
	brewPathMu.RLock()
	override := brewPathUser
	brewPathMu.RUnlock()
	if override != "" {
		path, err := exec.LookPath(override)
		if err != nil {
			return "", fmt.Errorf("configured brew path is invalid: %w", err)
		}
		return path, nil
	}

	brewPathOnce.Do(func() {
		if path, err := exec.LookPath("brew"); err == nil {
			brewPath = path
			return
		}

		candidates := []string{
			"/opt/homebrew/bin/brew",
			"/usr/local/bin/brew",
			"/home/linuxbrew/.linuxbrew/bin/brew",
		}
		for _, candidate := range candidates {
			if _, err := exec.LookPath(candidate); err == nil {
				brewPath = candidate
				return
			}
		}
		brewPathErr = errors.New("brew executable not found; tried PATH and common install paths")
	})

	if brewPathErr != nil {
		return "", brewPathErr
	}
	return brewPath, nil
}

func SetBrewPathOverride(path string) {
	brewPathMu.Lock()
	defer brewPathMu.Unlock()
	brewPathUser = strings.TrimSpace(path)
}

func SetCommandRuntimeOptions(maxConcurrency int, enableDebugLog bool) {
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}
	if maxConcurrency > 16 {
		maxConcurrency = 16
	}
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	if brewSem == nil || cap(brewSem) != maxConcurrency {
		brewSem = make(chan struct{}, maxConcurrency)
	}
	debugLog = enableDebugLog
}

func acquireBrewSlot(ctx context.Context) error {
	runtimeMu.Lock()
	if brewSem == nil {
		brewSem = make(chan struct{}, 4)
	}
	sem := brewSem
	runtimeMu.Unlock()
	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func releaseBrewSlot() {
	runtimeMu.Lock()
	sem := brewSem
	runtimeMu.Unlock()
	if sem == nil {
		return
	}
	select {
	case <-sem:
	default:
	}
}

func isDebugLogEnabled() bool {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	return debugLog
}

func writeDebugLog(line string) {
	if !isDebugLogEnabled() {
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	logDir := filepath.Join(homeDir, "Library", "Logs", "Gobrew")
	if mkErr := os.MkdirAll(logDir, 0755); mkErr != nil {
		return
	}
	logPath := filepath.Join(logDir, "gobrew-debug.log")
	file, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(time.Now().Format(time.RFC3339) + " " + line + "\n")
}

func runBrewCommand(ctx context.Context, args ...string) (string, string, error) {
	if err := acquireBrewSlot(ctx); err != nil {
		return "", err.Error(), err
	}
	defer releaseBrewSlot()
	brewExec, err := resolveBrewPath()
	if err != nil {
		return "", err.Error(), err
	}

	writeDebugLog(fmt.Sprintf("cmd: %s %s", brewExec, strings.Join(args, " ")))
	cmd := exec.CommandContext(ctx, brewExec, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		writeDebugLog(fmt.Sprintf("cmd err: %v", err))
	}
	return stdout.String(), stderr.String(), err
}

func runBrewCommandWithEvents(ctx context.Context, app *application.App, args ...string) (string, string, error) {
	if err := acquireBrewSlot(ctx); err != nil {
		return "", err.Error(), err
	}
	defer releaseBrewSlot()
	brewExec, err := resolveBrewPath()
	if err != nil {
		return "", err.Error(), err
	}

	writeDebugLog(fmt.Sprintf("cmd(events): %s %s", brewExec, strings.Join(args, " ")))
	cmd := exec.CommandContext(ctx, brewExec, args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("failed to start brew command: %w", err)
	}

	var stdoutBuf, stderrBuf strings.Builder
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go streamPipe(stdoutPipe, &stdoutBuf, app, stdoutDone)
	go streamPipe(stderrPipe, &stderrBuf, app, stderrDone)

	waitErr := cmd.Wait()
	<-stdoutDone
	<-stderrDone

	return stdoutBuf.String(), stderrBuf.String(), waitErr
}

func commandMessage(stdout, stderr string) string {
	msg := strings.TrimSpace(strings.TrimSpace(stdout) + "\n" + strings.TrimSpace(stderr))
	if msg == "" {
		return strings.TrimSpace(stderr)
	}
	return msg
}

func streamPipe(pipe interface{ Read([]byte) (int, error) }, dst *strings.Builder, app *application.App, done chan<- struct{}) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		dst.WriteString(line)
		dst.WriteString("\n")
		writeDebugLog(line)
		if app != nil {
			app.Event.Emit("brew-output", line)
		}
	}
	close(done)
}

func emitBrewComplete(app *application.App, success bool, details string, duration time.Duration) {
	if app == nil {
		return
	}

	trimmed := strings.TrimSpace(details)
	if success {
		if duration > 0 {
			app.Event.Emit("brew-complete", fmt.Sprintf("success (%s)", duration.String()))
			return
		}
		app.Event.Emit("brew-complete", "success")
		return
	}

	if trimmed == "" {
		trimmed = "command failed"
	}
	if duration > 0 {
		app.Event.Emit("brew-complete", fmt.Sprintf("failed (%s): %s", duration.String(), trimmed))
		return
	}
	app.Event.Emit("brew-complete", "failed: "+trimmed)
}
