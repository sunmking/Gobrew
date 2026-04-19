# Gobrew

[English](./README.md) | [中文](./docs/README.zh-CN.md)

## Download

- macOS installer and release assets: [GitHub Releases](https://github.com/sunmking/Gobrew/releases)

<p align="center">
  <img src="./docs/logo.png" alt="Gobrew Logo" width="160" />
</p>

<p align="center">
  <strong>Desktop command center for Homebrew.</strong><br/>
  Go + Wails + Vue, with clean workflows for install, update, maintain, and services.
</p>

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white">
  <img alt="Wails" src="https://img.shields.io/badge/Wails-v3-1E293B">
  <img alt="Vue" src="https://img.shields.io/badge/Vue-3-42B883?logo=vue.js&logoColor=white">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-22C55E">
</p>

## Why Gobrew

Homebrew CLI is fast and powerful.  
Gobrew makes it visual, batch-friendly, and easier to operate daily.

- Faster package search and installation
- Clear bulk update and cleanup workflows
- Unified taps, services, and Brewfile management
- Native Homebrew behavior with better UX

## Stack

- Backend: Go + Wails v3
- Frontend: Vue 3 + TypeScript + Pinia + Vue Router + Vue I18n
- Styling: Tailwind utilities + unified custom style system

## Quick Start

### For Users (Install and Use)

1. Open [GitHub Releases](https://github.com/sunmking/Gobrew/releases)
2. Download the latest `.dmg` installer
3. Install and launch Gobrew on macOS

### For Developers (Run Locally)

- macOS (primary target)
- Homebrew installed and available in `PATH`
- Go `1.25+`
- Node.js + npm
- Wails v3 CLI
- Task (`go-task`) for shortcut commands (recommended)

```bash
# install frontend deps
npm --prefix frontend install

# run desktop app
task dev
```

You can also run Wails directly:

```bash
wails3 dev -config ./build/config.yml -port 9245
```

## Core Workflows

### Explore
- search formulae/casks with categorized results
- inspect package details and dependencies
- bulk select and install

### Install
- inspect installed formulae/casks
- uninstall/reinstall packages
- monitor command output in real time

### Maintain
- manage taps (`add/remove/list/info`)
- run cleanup and autoremove related actions
- execute Brewfile workflows (`dump/list/check/cleanup/install/restore`)

### Services
- list and inspect running/background services
- bulk start/stop/restart operations
- observe streaming output and completion events

## Build and Release

```bash
task build
task package
```

## Command Cheatsheet

```bash
task dev            # desktop dev mode
task run            # run built desktop app
task build          # build desktop app
task package        # package distributables

task build:server   # build backend server mode
task run:server     # run backend server mode
task build:docker   # build docker image for server mode
task run:docker     # run docker image
```

## Developer Notes

### Project Layout

```text
Gobrew/
├── frontend/         # Vue application
├── services/         # Go service layer (brew, tap, service, bundle)
├── build/            # Wails config and packaging
├── docs/             # docs and visual assets
├── main.go           # app bootstrap and service registration
└── Taskfile.yml      # developer task entrypoints
```

### Runtime Model

- Frontend calls typed Wails bindings from service modules:
  - `BrewService`
  - `TapService`
  - `ServiceManager`
  - `BundleService`
- Long-running operations stream events to UI:
  - `brew-output`
  - `brew-complete`
- i18n is built-in (`en`, `zh`) and can be extended globally

## Roadmap

- virtualized rendering for very large search result sets
- deeper package insight panels and dependency graphs
- release automation and signed desktop artifacts
- broader platform/documentation coverage

## Contributing

PRs and issues are welcome.  
For larger features, open an issue first to align on scope.

## License

[MIT](./LICENSE)
