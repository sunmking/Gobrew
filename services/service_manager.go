package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ServiceManager struct {
	app *application.App
}

func NewServiceManager(app *application.App) *ServiceManager {
	return &ServiceManager{app: app}
}

func (s *ServiceManager) SetApp(app *application.App) {
	s.app = app
}

func (s *ServiceManager) List(ctx context.Context) ([]BrewServicesResult, error) {
	stdout, stderr, err := runBrewCommand(ctx, "services", "list")
	if err != nil && stdout == "" {
		return nil, &BrewError{Code: "SERVICES_LIST_FAILED", Message: "Failed to list services", Details: stderr}
	}

	results := make([]BrewServicesResult, 0)
	for _, line := range strings.Split(strings.TrimSpace(stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToLower(line), "name") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		exitCode := 0
		if len(fields) > 4 {
			if code, convErr := strconv.Atoi(fields[4]); convErr == nil {
				exitCode = code
			}
		}
		results = append(results, BrewServicesResult{
			Name:     fields[0],
			Status:   fields[1],
			User:     fields[2],
			File:     fields[3],
			ExitCode: exitCode,
		})
	}
	return results, nil
}

func (s *ServiceManager) Start(ctx context.Context, name string) error {
	return s.runServiceCommand(ctx, "SERVICE_START_FAILED", "Failed to start "+name, "services", "start", name)
}

func (s *ServiceManager) Stop(ctx context.Context, name string) error {
	return s.runServiceCommand(ctx, "SERVICE_STOP_FAILED", "Failed to stop "+name, "services", "stop", name)
}

func (s *ServiceManager) Restart(ctx context.Context, name string) error {
	return s.runServiceCommand(ctx, "SERVICE_RESTART_FAILED", "Failed to restart "+name, "services", "restart", name)
}

func (s *ServiceManager) StartAll(ctx context.Context) error {
	return s.runServiceCommand(ctx, "SERVICE_START_ALL_FAILED", "Failed to start all services", "services", "start", "--all")
}

func (s *ServiceManager) StopAll(ctx context.Context) error {
	return s.runServiceCommand(ctx, "SERVICE_STOP_ALL_FAILED", "Failed to stop all services", "services", "stop", "--all")
}

func (s *ServiceManager) RestartAll(ctx context.Context) error {
	return s.runServiceCommand(ctx, "SERVICE_RESTART_ALL_FAILED", "Failed to restart all services", "services", "restart", "--all")
}

func (s *ServiceManager) Cleanup(ctx context.Context) error {
	return s.runServiceCommand(ctx, "SERVICE_CLEANUP_FAILED", "Failed to cleanup services", "services", "cleanup")
}

func (s *ServiceManager) runServiceCommand(ctx context.Context, code, message string, args ...string) error {
	_, stderr, err := runBrewCommandWithEvents(ctx, s.app, args...)
	if err != nil {
		if s.app != nil {
			s.app.Event.Emit("brew-complete", fmt.Sprintf(`{"success":false,"error":"%s"}`, stderr))
		}
		return &BrewError{Code: code, Message: message, Details: stderr}
	}
	if s.app != nil {
		s.app.Event.Emit("brew-complete", `{"success":true}`)
	}
	return nil
}
