import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { Get, Import, Reset, Save, ValidateBrewPath } from '../../bindings/changeme/services/configservice.js'
import { AppConfig, BrewPathValidation } from '../../bindings/changeme/services/models.js'

export type ThemeMode = 'dark' | 'light'
export type AutoUpdateInterval = 'off' | '1h' | '6h' | '12h' | '24h'

const DEFAULT_LANG = navigator.language.startsWith('zh') ? 'zh' : 'en'

function applyTheme(theme: ThemeMode) {
  const root = document.documentElement
  root.dataset.theme = theme
  root.classList.toggle('dark', theme === 'dark')
}

export const useSettingsStore = defineStore('settings', () => {
  const initialized = ref(false)
  const loading = ref(false)
  const saving = ref(false)
  const config = ref(new AppConfig({
    language: DEFAULT_LANG,
    theme: 'dark',
    brew_path: '',
    brewfile_path: '',
    auto_update_interval: 'off',
    log_max_lines: 1000,
    check_updates_on_launch: true,
  }))

  const language = computed({
    get: () => config.value.language || DEFAULT_LANG,
    set: (value: string) => { config.value.language = value },
  })
  const brewPath = computed({
    get: () => config.value.brew_path || '',
    set: (value: string) => { config.value.brew_path = value },
  })
  const brewFilePath = computed({
    get: () => config.value.brewfile_path || '',
    set: (value: string) => { config.value.brewfile_path = value },
  })
  const autoUpdate = computed({
    get: () => (config.value.auto_update_interval || 'off') as AutoUpdateInterval,
    set: (value: AutoUpdateInterval) => { config.value.auto_update_interval = value },
  })
  const theme = computed({
    get: () => (config.value.theme || 'dark') as ThemeMode,
    set: (value: ThemeMode) => { config.value.theme = value },
  })
  const logMaxLines = computed({
    get: () => Math.max(100, Number(config.value.log_max_lines || 1000)),
    set: (value: number) => { config.value.log_max_lines = Math.max(100, value) },
  })
  const checkUpdatesOnLaunch = computed({
    get: () => Boolean(config.value.check_updates_on_launch),
    set: (value: boolean) => { config.value.check_updates_on_launch = value },
  })

  async function load() {
    if (loading.value) return
    loading.value = true
    try {
      config.value = await Get()
      applyTheme(theme.value)
      initialized.value = true
    } finally {
      loading.value = false
    }
  }

  async function patchAndSave(patch: Partial<AppConfig>) {
    saving.value = true
    try {
      const next = new AppConfig({
        ...config.value,
        ...patch,
      })
      config.value = await Save(next)
      applyTheme(theme.value)
    } finally {
      saving.value = false
    }
  }

  function setLanguage(value: string) {
    return patchAndSave({ language: value })
  }

  function setTheme(value: ThemeMode) {
    return patchAndSave({ theme: value })
  }

  function setBrewPath(value: string) {
    return patchAndSave({ brew_path: value.trim() })
  }

  function setBrewFilePath(value: string) {
    return patchAndSave({ brewfile_path: value.trim() })
  }

  function setAutoUpdate(value: AutoUpdateInterval) {
    return patchAndSave({ auto_update_interval: value })
  }

  function setLogMaxLines(value: number) {
    return patchAndSave({ log_max_lines: Math.max(100, value) })
  }

  function setCheckUpdatesOnLaunch(value: boolean) {
    return patchAndSave({ check_updates_on_launch: value })
  }

  async function resetAll() {
    config.value = await Reset()
    applyTheme(theme.value)
  }

  async function importFromContent(content: string) {
    config.value = await Import(content)
    applyTheme(theme.value)
  }

  async function validateBrewPath(path: string): Promise<BrewPathValidation> {
    return ValidateBrewPath(path)
  }

  return {
    initialized,
    loading,
    saving,
    config,
    language,
    brewPath,
    brewFilePath,
    autoUpdate,
    theme,
    logMaxLines,
    checkUpdatesOnLaunch,
    load,
    setLanguage,
    setTheme,
    setBrewPath,
    setBrewFilePath,
    setAutoUpdate,
    setLogMaxLines,
    setCheckUpdatesOnLaunch,
    resetAll,
    importFromContent,
    validateBrewPath,
  }
})
