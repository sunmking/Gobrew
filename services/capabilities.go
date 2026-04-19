package services

import (
	"context"
	"strings"
	"sync"
	"time"
)

type BrewCapabilities struct {
	BrewVersion                 string    `json:"brew_version"`
	SupportsListJSONV2          bool      `json:"supports_list_json_v2"`
	SupportsTapJSON             bool      `json:"supports_tap_json"`
	SupportsOutdatedJSONV2      bool      `json:"supports_outdated_json_v2"`
	SupportsInfoInstalledJSONV2 bool      `json:"supports_info_installed_json_v2"`
	DetectedAt                  time.Time `json:"detected_at"`
}

var (
	capabilitiesOnce sync.Once
	capabilities     BrewCapabilities
)

func getBrewCapabilities(ctx context.Context) BrewCapabilities {
	capabilitiesOnce.Do(func() {
		capabilities = detectBrewCapabilities(ctx)
	})
	return capabilities
}

func detectBrewCapabilities(ctx context.Context) BrewCapabilities {
	caps := BrewCapabilities{
		DetectedAt: time.Now(),
	}

	versionOut, _, _ := runBrewCommand(ctx, "--version")
	if firstLine := strings.TrimSpace(strings.Split(strings.TrimSpace(versionOut), "\n")[0]); firstLine != "" {
		caps.BrewVersion = firstLine
	}

	listOut, listErrText, listErr := runBrewCommand(ctx, "list", "--formula", "--json=v2")
	caps.SupportsListJSONV2 = listErr == nil || !containsUnsupportedJSONOption(commandMessage(listOut, listErrText))

	tapOut, tapErrText, tapErr := runBrewCommand(ctx, "tap", "--json")
	caps.SupportsTapJSON = tapErr == nil || !containsUnsupportedJSONOption(commandMessage(tapOut, tapErrText))

	outdatedOut, outdatedErrText, outdatedErr := runBrewCommand(ctx, "outdated", "--json=v2")
	caps.SupportsOutdatedJSONV2 = outdatedErr == nil || !containsUnsupportedJSONOption(commandMessage(outdatedOut, outdatedErrText))

	// Keep capability for compatibility, but avoid expensive detection commands here.
	caps.SupportsInfoInstalledJSONV2 = caps.SupportsListJSONV2

	return caps
}
