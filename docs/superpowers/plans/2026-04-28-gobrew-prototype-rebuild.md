# Gobrew Prototype Rebuild Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild Gobrew around the high-fidelity prototype with five functional entries: All Packages, Update & Cleanup, Taps, Services, and Brewfile.

**Architecture:** Preserve the existing Wails service boundaries and add only narrow backend support for cask package info and Brewfile writing. Replace the current frontend information architecture with a compact prototype-aligned app shell, shared operational UI components, and five routed pages connected to existing Pinia stores and generated Wails bindings.

**Tech Stack:** Go, Wails 3, Vue 3, Pinia, Vue Router, TypeScript, Tailwind CSS, lucide-vue-next.

---

## File Structure

- Modify `services/models.go`: add unified package detail structs for formula and cask metadata.
- Modify `services/brew_service.go`: add `PackageInfo(ctx, name, packageType)` and cask JSON parsing.
- Modify `services/bundle_service.go`: add `WriteBrewfile(filePath, content)`.
- Add/modify tests in `services/brew_service_test.go` and `services/bundle_service_test.go`.
- Regenerate `frontend/bindings/changeme/services/*`.
- Modify `frontend/src/router/index.ts`: route to the five prototype pages and package detail.
- Modify `frontend/src/App.vue`: use the prototype shell layout.
- Modify `frontend/src/style.css`: replace global tokens and shared component classes with prototype-aligned styling.
- Modify `frontend/src/components/layout/Sidebar.vue`: render five-entry prototype navigation.
- Create `frontend/src/components/layout/TitleBar.vue`.
- Create shared UI components under `frontend/src/components/common/`: `BrewButton.vue`, `SegmentedControl.vue`, `StatCard.vue`, `StatusPill.vue`, `TerminalPanel.vue`, `ToolbarSearch.vue`.
- Create package components under `frontend/src/components/packages/`: `PackageTable.vue`.
- Modify stores under `frontend/src/stores/` only where needed to expose package info, Brewfile write, tap details, and action states.
- Replace page implementations in `frontend/src/pages/`: `AllPackagesPage.vue`, `PackageDetailPage.vue`, `UpdateCleanupPage.vue`, `TapsPage.vue`, `ServicesPage.vue`, `BrewfilePage.vue`.

## Task 1: Backend Package Info And Brewfile Write

**Files:**
- Modify: `services/models.go`
- Modify: `services/brew_service.go`
- Modify: `services/bundle_service.go`
- Test: `services/brew_service_test.go`
- Test: `services/bundle_service_test.go`

- [ ] **Step 1: Write failing backend tests**

Add tests that call pure parsing helpers and path-safe Brewfile writing behavior:

```go
func TestParsePackageInfoJSONIncludesCask(t *testing.T) {
	raw := `{"formulae":[],"casks":[{"name":"visual-studio-code","full_name":"visual-studio-code","tap":"homebrew/cask","desc":"Open-source code editor","homepage":"https://code.visualstudio.com/","version":"1.99.0","installed":"1.98.2","auto_updates":true,"token":"visual-studio-code"}]}`
	info, err := parsePackageInfoJSON(raw, "visual-studio-code", "cask")
	if err != nil {
		t.Fatalf("parsePackageInfoJSON returned error: %v", err)
	}
	if info.Type != "cask" || info.Name != "visual-studio-code" || info.CurrentVersion != "1.99.0" || info.InstalledVersion != "1.98.2" {
		t.Fatalf("unexpected cask info: %+v", info)
	}
	if !info.AutoUpdates {
		t.Fatalf("expected auto_updates to be true")
	}
}

func TestParsePackageInfoJSONIncludesFormula(t *testing.T) {
	raw := `{"formulae":[{"name":"wget","full_name":"wget","tap":"homebrew/core","desc":"Internet file retriever","homepage":"https://www.gnu.org/software/wget/","license":"GPL-3.0","linked_keg":"1.24.5","pinned":false,"installed":[{"version":"1.24.5","installed_as_dependency":false,"installed_on_request":true}],"versions":{"stable":"1.24.5","head":"HEAD"},"dependencies":["openssl@3","zlib"]}],"casks":[]}`
	info, err := parsePackageInfoJSON(raw, "wget", "formula")
	if err != nil {
		t.Fatalf("parsePackageInfoJSON returned error: %v", err)
	}
	if info.Type != "formula" || info.Name != "wget" || info.CurrentVersion != "1.24.5" || info.InstalledVersion != "1.24.5" {
		t.Fatalf("unexpected formula info: %+v", info)
	}
	if len(info.Dependencies) != 2 || info.Dependencies[0] != "openssl@3" {
		t.Fatalf("unexpected dependencies: %+v", info.Dependencies)
	}
}

func TestWriteBrewfileWritesResolvedPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Brewfile")
	service := NewBundleService(nil)
	if err := service.WriteBrewfile(path, "brew \"wget\"\n"); err != nil {
		t.Fatalf("WriteBrewfile returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "brew \"wget\"\n" {
		t.Fatalf("unexpected Brewfile content: %q", string(data))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./services -run 'TestParsePackageInfoJSONIncludes|TestWriteBrewfileWritesResolvedPath'`

Expected: FAIL because `parsePackageInfoJSON` and `WriteBrewfile` do not exist.

- [ ] **Step 3: Implement backend additions**

Add a `PackageInfoResult` DTO to `services/models.go` with package type, name, full name, tap, description, homepage, license, current version, installed version, linked keg, pinned, auto updates, token, installed kegs, and dependencies.

Add `BrewService.PackageInfo(ctx, name, packageType)` in `services/brew_service.go`. It should run `brew info --json=v2` with `--cask` when `packageType == "cask"` and parse the result through `parsePackageInfoJSON`.

Add `BundleService.WriteBrewfile(filePath, content)` in `services/bundle_service.go`. It should use `resolveBrewfilePath`, create the parent directory with `0755`, and write content with `0644`.

- [ ] **Step 4: Run backend focused tests**

Run: `go test ./services -run 'TestParsePackageInfoJSONIncludes|TestWriteBrewfileWritesResolvedPath'`

Expected: PASS.

- [ ] **Step 5: Run backend package tests**

Run: `go test ./services`

Expected: PASS.

## Task 2: Regenerate Wails Bindings

**Files:**
- Modify: `frontend/bindings/changeme/services/brewservice.ts`
- Modify: `frontend/bindings/changeme/services/bundleservice.ts`
- Modify: `frontend/bindings/changeme/services/models.ts`

- [ ] **Step 1: Generate bindings**

Run: `wails3 generate bindings`

Expected: generated TypeScript bindings include `BrewService.PackageInfo`, `BundleService.WriteBrewfile`, and `PackageInfoResult`.

- [ ] **Step 2: Inspect generated bindings**

Run: `rg -n "PackageInfo|WriteBrewfile|PackageInfoResult" frontend/bindings/changeme/services`

Expected: matching generated method and type names appear in service bindings.

## Task 3: Frontend Shell, Routes, And Shared Components

**Files:**
- Modify: `frontend/src/App.vue`
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/style.css`
- Modify: `frontend/src/components/layout/Sidebar.vue`
- Create: `frontend/src/components/layout/TitleBar.vue`
- Create: `frontend/src/components/common/BrewButton.vue`
- Create: `frontend/src/components/common/SegmentedControl.vue`
- Create: `frontend/src/components/common/StatCard.vue`
- Create: `frontend/src/components/common/StatusPill.vue`
- Create: `frontend/src/components/common/TerminalPanel.vue`
- Create: `frontend/src/components/common/ToolbarSearch.vue`

- [ ] **Step 1: Implement routes**

Replace the router table with:

```ts
const routes = [
  { path: '/', name: 'packages', component: () => import('@/pages/AllPackagesPage.vue') },
  { path: '/update-cleanup', name: 'update-cleanup', component: () => import('@/pages/UpdateCleanupPage.vue') },
  { path: '/taps', name: 'taps', component: () => import('@/pages/TapsPage.vue') },
  { path: '/services', name: 'services', component: () => import('@/pages/ServicesPage.vue') },
  { path: '/brewfile', name: 'brewfile', component: () => import('@/pages/BrewfilePage.vue') },
  { path: '/packages/:type/:name', name: 'package-detail', component: () => import('@/pages/PackageDetailPage.vue') },
  { path: '/install', redirect: '/' },
  { path: '/installed', redirect: '/' },
  { path: '/explore', redirect: '/' },
  { path: '/upgrade', redirect: '/update-cleanup' },
  { path: '/update', redirect: '/update-cleanup' },
  { path: '/maintain', redirect: '/update-cleanup' },
  { path: '/cleanup', redirect: '/update-cleanup' },
  { path: '/bundle', redirect: '/brewfile' },
  { path: '/settings', redirect: '/' },
]
```

- [ ] **Step 2: Implement shell components**

`App.vue` should render `TitleBar`, `Sidebar`, `router-view`, and `OperationStatusBar` inside `.window-shell`, `.app-body`, and `.content` wrappers.

`TitleBar.vue` should render the Gobrew title, traffic-light dots, and a theme toggle that switches `document.documentElement.dataset.theme` between `light` and `dark`.

`Sidebar.vue` should render five entries with lucide icons and active state based on route path.

- [ ] **Step 3: Implement shared components**

Create small components for prototype controls. Each component should accept simple props and emit native-style events:

- `BrewButton`: `variant`, `disabled`, default slot.
- `SegmentedControl`: `modelValue`, `options`, emits `update:modelValue`.
- `StatCard`: `label`, `value`, `sub`, `tone`.
- `StatusPill`: `status`, default slot.
- `TerminalPanel`: `lines`.
- `ToolbarSearch`: `modelValue`, `placeholder`, emits `update:modelValue` and `submit`.

- [ ] **Step 4: Replace global styling**

Move `frontend/src/style.css` toward the prototype tokens and classes while keeping Tailwind directives. Define `.window-shell`, `.titlebar`, `.app-body`, `.app-sidebar`, `.content`, `.toolbar`, `.pkg-table`, `.pill`, `.stat-grid`, `.stat-card`, `.terminal`, and shared button styles.

- [ ] **Step 5: Run frontend build**

Run: `npm run build` from `frontend`.

Expected: build may fail because pages are not replaced yet; route/component type errors should be fixed before moving to page tasks.

## Task 4: All Packages And Detail Pages

**Files:**
- Create: `frontend/src/pages/AllPackagesPage.vue`
- Create: `frontend/src/pages/PackageDetailPage.vue`
- Create/modify: `frontend/src/components/packages/PackageTable.vue`
- Modify: `frontend/src/stores/search.ts`
- Modify: `frontend/src/types/brew.ts`

- [ ] **Step 1: Extend frontend types and search store**

Add `PackageInfoResult` and a package row view model type to `frontend/src/types/brew.ts`. Add `packageInfo(name, type)` action or direct binding call usage for detail.

- [ ] **Step 2: Implement `PackageTable.vue`**

The table should render package type, name, description, installed version, latest version, status pill, and row actions. It should emit `select`, `install`, `upgrade`, and `uninstall`.

- [ ] **Step 3: Implement `AllPackagesPage.vue`**

Load installed and outdated stores on mount. Search via `useSearchStore`. Merge formulae/casks/outdated/search into rows keyed by `type:name`. Support segmented filters `all`, `installed`, `updates`. Navigate to `/packages/:type/:name` on row select.

- [ ] **Step 4: Implement `PackageDetailPage.vue`**

Load `BrewService.PackageInfo(name, type)`, `Deps(name, false)` for formulae, and `Uses(name, true, false)` for formulae. Render prototype-like header, metadata cards, dependencies, dependents, and terminal panel using the log store.

- [ ] **Step 5: Run frontend build**

Run: `npm run build` from `frontend`.

Expected: PASS after type and component fixes.

## Task 5: Update, Taps, Services, And Brewfile Pages

**Files:**
- Create: `frontend/src/pages/UpdateCleanupPage.vue`
- Modify: `frontend/src/pages/TapsPage.vue`
- Modify: `frontend/src/pages/ServicesPage.vue`
- Create: `frontend/src/pages/BrewfilePage.vue`
- Modify: `frontend/src/stores/bundle.ts`
- Modify: `frontend/src/stores/taps.ts`

- [ ] **Step 1: Implement `UpdateCleanupPage.vue`**

Use outdated data, `CleanupPreview`, `AutoRemovePreview`, `UpgradeAll`, `Cleanup`, and `AutoRemove`. Render stats, tabs, tables/output panels, and action buttons.

- [ ] **Step 2: Implement `TapsPage.vue`**

Use tap list and selected tap detail. Render search, add tap action, tap cards/list rows, lazy details with formula/cask counts, and remove actions.

- [ ] **Step 3: Implement `ServicesPage.vue`**

Use services store. Render service table with status pills and row actions for start, stop, restart plus toolbar refresh/restart all.

- [ ] **Step 4: Implement `BrewfilePage.vue` and bundle store write**

Use `ReadBrewfile`, `WriteBrewfile`, `Check`, `Dump`, `Restore`, and `CleanupPreview`. Render editor, action stack, missing dependency preview, cleanup preview, and stats.

- [ ] **Step 5: Run frontend build**

Run: `npm run build` from `frontend`.

Expected: PASS.

## Task 6: Final Verification

**Files:**
- All changed frontend and backend files.

- [ ] **Step 1: Format Go code**

Run: `gofmt -w services/models.go services/brew_service.go services/bundle_service.go services/brew_service_test.go services/bundle_service_test.go`

Expected: command exits 0.

- [ ] **Step 2: Run frontend build**

Run: `npm run build` from `frontend`.

Expected: PASS.

- [ ] **Step 3: Run Go tests**

Run: `go test ./...`

Expected: PASS. macOS linker warnings are acceptable if tests pass.

- [ ] **Step 4: Check git status**

Run: `git status --short`

Expected: only intentional source, binding, test, and plan changes appear.
