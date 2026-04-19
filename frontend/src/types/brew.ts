export interface FormulaInfo {
  name: string
  full_name: string
  tap: string
  desc: string
  homepage: string
  license: string
  versions: {
    stable: string
    head: string
  }
  installed: Array<{
    version: string
    installed_as_dependency: boolean
    installed_on_request: boolean
  }>
  linked_keg: string | null
  pinned: boolean
  dependencies: string[]
  build_dependencies: string[]
}

export interface FormulaInstalled {
  name: string
  full_name: string
  tap: string
  desc: string
  homepage: string
  license: string
  versions: { stable: string }
  installed: Array<{
    version: string
    installed_as_dependency: boolean
    installed_on_request: boolean
  }>
  linked_keg: string | null
  pinned: boolean
  dependencies: string[]
}

export interface CaskInstalled {
  name: string
  full_name: string
  tap: string
  desc: string
  homepage: string
  version: string
  installed: string
  auto_updates: boolean
  token: string
}

export interface OutdatedFormula {
  name: string
  installed_versions: string[]
  current_version: string
  pinned: boolean
  pinned_version: string | null
}

export interface OutdatedCask {
  name: string
  installed_version: string
  current_version: string
}

export interface BrewService {
  name: string
  status: 'started' | 'stopped' | 'error' | 'none'
  user: string
  file: string
  exit_code: number
}

export interface TapInfo {
  name: string
  remote: string
  custom_remote: boolean
}

export interface TapDetail {
  name: string
  remote: string
  custom_remote: boolean
  formula_names: string[]
  cask_tokens: string[]
  last_commit: string
  branch: string
}

export interface SearchResult {
  formulae: Array<{
    name: string
    full_name: string
    desc: string
    tap: string
  }>
  casks: Array<{
    name: string
    full_name: string
    desc: string
    tap: string
  }>
}

export interface BrewError {
  code: string
  message: string
  details?: string
}

export interface LogLine {
  id: number
  text: string
  type: 'stdout' | 'stderr' | 'system'
  timestamp: number
}
