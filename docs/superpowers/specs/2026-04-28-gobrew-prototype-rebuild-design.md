# Gobrew Prototype Rebuild Design

Date: 2026-04-28

## Decision

Rebuild Gobrew around the high-fidelity prototype in `docs/gobrew-prototype.html` using approach 1: prototype-led UI and interaction refactor, while preserving existing Wails service boundaries and adding only the backend aggregation methods needed by the new experience.

The primary navigation will be replaced with the prototype's five entries:

- All Packages
- Update & Cleanup
- Taps
- Services
- Brewfile

The existing Dashboard, Install, Explore, Maintain, and Settings navigation model will no longer be the primary information architecture. Install and search workflows move into All Packages.

## Goals

- Match the prototype's desktop application shell: compact title bar, left sidebar, dense table-first content, status pills, toolbar search, segmented filters, terminal output, and light/dark tokens.
- Make the rebuilt app functional, not static. Existing brew, tap, services, and bundle commands remain connected to the UI.
- Keep the implementation conservative by extending current Go services instead of replacing the backend architecture.
- Prefer dense operational screens over marketing-style or card-heavy layouts.
- Preserve command output streaming through existing `brew-output` and `brew-complete` events.

## Non-Goals

- Do not redesign the backend into a new service architecture.
- Do not implement a full package registry cache or background indexer.
- Do not add cloud accounts, sync, notifications, or advanced Homebrew analytics.
- Do not keep the current Dashboard/Install/Explore/Maintain/Settings layout as parallel navigation.

## Current Context

Gobrew is a Wails 3 application with a Vue 3, Pinia, Vue Router, Tailwind, and lucide frontend. The backend already exposes useful service boundaries:

- `BrewService`: installed packages, search, info, install, uninstall, upgrade, update, deps, uses, outdated, cleanup, autoremove, doctor.
- `TapService`: list, add, remove, details, and package actions.
- `BundleService`: dump, restore, check, cleanup preview, cleanup, list, read Brewfile.
- `ServiceManager`: list, start, stop, restart, bulk service actions, cleanup.

The prototype is a single HTML file that defines the target interaction model and visual system, including the app shell, package table, package detail screen, update/cleanup screen, taps list, services table, and Brewfile editor.

## Architecture

### Application Shell

The Vue app will use a single shell layout:

- `AppShell` owns the title bar, sidebar, routed content area, and operation status bar.
- `TitleBar` renders traffic-light affordances, centered Gobrew title, theme toggle, and action icons that map to supported app actions.
- `Sidebar` renders only the five prototype entries, grouped as package management and resources.
- `router-view` renders one of the five feature pages.
- `OperationStatusBar` remains available for long-running brew command feedback.

### Routes

The route model becomes:

- `/` -> All Packages
- `/update-cleanup` -> Update & Cleanup
- `/taps` -> Taps
- `/services` -> Services
- `/brewfile` -> Brewfile
- `/packages/:type/:name` -> Package Detail

`type` is `formula` or `cask`. Legacy routes can redirect to the closest new route during the transition.

### Shared Frontend Components

The refactor will extract prototype-aligned components:

- `ToolbarSearch`: search input with icon and compact sizing.
- `SegmentedControl`: filters such as all, installed, and updates.
- `PackageTable`: dense package list for formulae and casks.
- `StatusPill`: installed, update available, not installed, running, stopped, and error states.
- `StatCard`: compact metric cards for update and cleanup summaries.
- `TerminalPanel`: streamed command output with prompt/output/success/error styling.
- `DetailPanel`: reusable side-panel card for metadata blocks.
- `BrewButton`: shared button styling for primary, secondary, ghost, and danger actions.

Page components should stay thin and delegate display details to these shared components.

## Data Flow

### All Packages

All Packages combines installed formulae, installed casks, outdated formulae, outdated casks, and search results into one table model. The page will:

- Load installed and outdated data on mount.
- Search Homebrew when the user enters a query.
- Merge package state by name and type.
- Show installed version, latest/current version, description, tap, and status.
- Support filters for all, installed, and update available.
- Navigate to package detail when a row is selected.
- Trigger install, uninstall, or upgrade actions from row actions where appropriate.

The backend should expose a small aggregation DTO if the frontend merge logic becomes too complex. The initial implementation may merge in Pinia if the existing stores already provide the needed data.

### Package Detail

Package Detail loads package metadata using `brew info --json=v2`. It also loads dependencies and installed dependents using the existing `Deps` and `Uses` methods where useful.

Formula detail will use the existing `BrewService.Info` result. Cask detail needs a comparable backend path because the current `Info` method only maps formulae. The cask info response should include name, token, tap, description, homepage, version, installed version, and auto-update status.

### Update & Cleanup

Update & Cleanup combines:

- `Outdated` for update counts and update table.
- `CleanupPreview` for cleanup candidate output.
- `AutoRemovePreview` for removable dependency output.
- `Upgrade`, `UpgradeAll`, `Cleanup`, and `AutoRemove` for actions.

The UI has three tabs: pending updates, cleanup, and cache/dependencies. If Homebrew cannot provide exact reclaim size, the UI will show item counts and command output rather than inventing a size.

### Taps

Taps uses `TapService.List` and `TapService.Details`. The list page shows name, remote URL, official/custom state, and package count when available from details. Add and remove actions continue using existing methods.

To avoid expensive detail calls for every tap on initial load, the first version can show counts after selecting a tap or after a background enrichment pass.

### Services

Services uses `ServiceManager.List` and service actions. The table shows name, status, user, file, exit code, and row actions for start, stop, and restart. Bulk restart and refresh follow the prototype toolbar.

### Brewfile

Brewfile uses `BundleService.ReadBrewfile`, `Check`, `Dump`, `Restore`, `CleanupPreview`, and `Cleanup`. The editor reads the default `~/Brewfile` unless the user chooses another path later.

The first difference preview will be based on `brew bundle check --verbose` missing dependencies and cleanup preview output. It will not attempt a perfect semantic parser for every Brewfile directive.

The backend should add a write method so the editor can save Brewfile contents through Wails rather than relying only on dump.

## Backend Changes

Backend changes should be additive:

- Add cask support to package info, either by extending `BrewService.Info` or adding a `PackageInfo(ctx, name, type)` method.
- Add `BundleService.WriteBrewfile(filePath, content)` with the same path resolution rules as read.
- Add DTOs only where they reduce frontend merging complexity.
- Keep command execution in `command_runner.go` unchanged unless tests expose a streaming or error-handling issue.
- Preserve existing JSON fallback behavior for older Homebrew versions.

## State Management

Pinia stores should align with the five domains:

- `packages`: installed, outdated, search, detail, selected filters.
- `maintenance`: cleanup, autoremove, and update actions.
- `taps`: tap list and selected tap details.
- `services`: service rows and action state.
- `bundle`: Brewfile content, check result, cleanup preview, and restore/dump actions.

Existing stores may be renamed or folded into these domains when it makes the new pages simpler. Avoid broad rewrites of business logic that do not serve the prototype rebuild.

## Error Handling

- Backend methods continue returning `BrewError` with code, message, and details.
- Pages show concise toast errors and keep command details in `TerminalPanel` or status output.
- Empty states distinguish between "no results", "not loaded", and "command failed".
- Long-running operations disable only the relevant buttons, not the entire app.
- If Homebrew is missing, the app should show a recoverable state rather than a blank page.

## Styling

Global styling will move toward the prototype tokens:

- `--bg`, `--surface`, `--fg`, `--muted`, `--border`, `--accent`, `--accent-bg`, `--danger`, `--warning`, `--success`
- Compact radii: 4px, 6px, and 10px
- Dense table typography with system UI and monospace package names
- Light and dark theme support through a root theme attribute or class

Use lucide icons in Vue components instead of copying inline SVG from the prototype.

## Testing And Verification

Backend:

- Unit tests for any new parsing, package info mapping, and Brewfile write path resolution.
- Run `go test ./...`.

Frontend:

- Run `npm run build` from `frontend`.
- Verify TypeScript bindings if backend DTOs or methods change.

Bindings:

- Run `wails3 generate bindings` after backend service API changes.

Manual verification:

- All five sidebar entries render.
- Package search, installed filter, update filter, and detail navigation work.
- Update, cleanup, tap, service, and Brewfile actions emit visible command feedback.
- Light and dark themes are readable.

## Open Implementation Notes

- Some legacy routes can stay as redirects for compatibility, but they should not appear in navigation.
- Exact cleanup size is optional because Homebrew dry-run output may not provide reliable structured size data.
- Tap package counts can be loaded lazily to avoid slow startup.
- Brewfile semantic diff remains intentionally shallow in the first rebuild.
