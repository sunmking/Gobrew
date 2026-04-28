# Changelog

All notable changes to this project are documented in this file.

## v1.2.0 - 2026-04-29

### Added
- Added backend `ConfigService` with persistent application config.
- Added config import, export, and reset capabilities.
- Added brew path validation and runtime brew path override support.
- Added Homebrew version display in sidebar footer.
- Added tests for config service behavior.

### Changed
- Fully redesigned settings feature into a new Settings Center.
- Migrated frontend settings state to backend-driven single source of truth.
- Applied config-driven theme and language initialization at app startup.
- Connected Brewfile default path and log max lines to global config.
- Completed broad i18n integration across active pages and common UI.
- Updated notification button behavior to trigger update checks with status feedback.

### Removed
- Removed legacy pages not used by the active routed experience:
  - `BundlePage`
  - `CleanupPage`
  - `ExplorePage`
  - `HomePage`
  - `InstallPage`
  - `InstalledPage`
  - `MaintainPage`
  - `SystemPage`
  - `UpdatePage`

### Technical
- Regenerated Wails bindings for the new config service APIs.
- Refined package/search/install flow and waterfall loading continuity from prior iteration work.

