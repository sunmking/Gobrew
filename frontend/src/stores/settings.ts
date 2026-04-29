import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { Get, Import, Reset, Save, ValidateBrewPath } from '../../bindings/changeme/services/configservice.js'
import { AppConfig, BrewPathValidation } from '../../bindings/changeme/services/models.js'

export type ThemeMode = 'dark' | 'light'
export type AutoUpdateInterval = 'off' | '1h' | '6h' | '12h' | '24h'
export type RowDensity = 'compact' | 'default' | 'comfortable'
export type ToastPosition = 'tl' | 'tr' | 'bl' | 'br'

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
    launch_at_login: false,
    show_in_menu_bar: true,
    restore_last_page: true,
    last_page_path: '/',
    backup_before_update: true,
    cleanup_after_update: false,
    notify_updates: true,
    notify_operations: true,
    notify_services: true,
    notify_errors: true,
    show_toasts: true,
    toast_position: 'br',
    keep_notification_history: true,
    toast_duration_ms: 2500,
    max_concurrency: 4,
    debug_log: false,
    accent_hue: 255,
    ui_font_size: 13,
    row_density: 'default',
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
  const launchAtLogin = computed({
    get: () => Boolean(config.value.launch_at_login),
    set: (value: boolean) => { config.value.launch_at_login = value },
  })
  const showInMenuBar = computed({
    get: () => Boolean(config.value.show_in_menu_bar),
    set: (value: boolean) => { config.value.show_in_menu_bar = value },
  })
  const restoreLastPage = computed({
    get: () => Boolean(config.value.restore_last_page),
    set: (value: boolean) => { config.value.restore_last_page = value },
  })
  const lastPagePath = computed({
    get: () => config.value.last_page_path || '/',
    set: (value: string) => { config.value.last_page_path = value || '/' },
  })
  const backupBeforeUpdate = computed({
    get: () => Boolean(config.value.backup_before_update),
    set: (value: boolean) => { config.value.backup_before_update = value },
  })
  const cleanupAfterUpdate = computed({
    get: () => Boolean(config.value.cleanup_after_update),
    set: (value: boolean) => { config.value.cleanup_after_update = value },
  })
  const notifyUpdates = computed({
    get: () => Boolean(config.value.notify_updates),
    set: (value: boolean) => { config.value.notify_updates = value },
  })
  const notifyOperations = computed({
    get: () => Boolean(config.value.notify_operations),
    set: (value: boolean) => { config.value.notify_operations = value },
  })
  const notifyServices = computed({
    get: () => Boolean(config.value.notify_services),
    set: (value: boolean) => { config.value.notify_services = value },
  })
  const notifyErrors = computed({
    get: () => Boolean(config.value.notify_errors),
    set: (value: boolean) => { config.value.notify_errors = value },
  })
  const showToasts = computed({
    get: () => Boolean(config.value.show_toasts),
    set: (value: boolean) => { config.value.show_toasts = value },
  })
  const toastPosition = computed({
    get: () => (config.value.toast_position || 'br') as ToastPosition,
    set: (value: ToastPosition) => { config.value.toast_position = value },
  })
  const keepNotificationHistory = computed({
    get: () => Boolean(config.value.keep_notification_history),
    set: (value: boolean) => { config.value.keep_notification_history = value },
  })
  const toastDurationMs = computed({
    get: () => Math.max(1000, Number(config.value.toast_duration_ms || 2500)),
    set: (value: number) => { config.value.toast_duration_ms = Math.max(1000, value) },
  })
  const maxConcurrency = computed({
    get: () => Math.max(1, Number(config.value.max_concurrency || 4)),
    set: (value: number) => { config.value.max_concurrency = Math.max(1, value) },
  })
  const debugLog = computed({
    get: () => Boolean(config.value.debug_log),
    set: (value: boolean) => { config.value.debug_log = value },
  })
  const accentHue = computed({
    get: () => Number(config.value.accent_hue || 255),
    set: (value: number) => { config.value.accent_hue = value },
  })
  const uiFontSize = computed({
    get: () => Number(config.value.ui_font_size || 13),
    set: (value: number) => { config.value.ui_font_size = value },
  })
  const rowDensity = computed({
    get: () => (config.value.row_density || 'default') as RowDensity,
    set: (value: RowDensity) => { config.value.row_density = value },
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
  function setLaunchAtLogin(value: boolean) {
    return patchAndSave({ launch_at_login: value })
  }
  function setShowInMenuBar(value: boolean) {
    return patchAndSave({ show_in_menu_bar: value })
  }
  function setRestoreLastPage(value: boolean) {
    return patchAndSave({ restore_last_page: value })
  }
  function setLastPagePath(value: string) {
    return patchAndSave({ last_page_path: value || '/' })
  }
  function setBackupBeforeUpdate(value: boolean) {
    return patchAndSave({ backup_before_update: value })
  }
  function setCleanupAfterUpdate(value: boolean) {
    return patchAndSave({ cleanup_after_update: value })
  }
  function setNotifyUpdates(value: boolean) {
    return patchAndSave({ notify_updates: value })
  }
  function setNotifyOperations(value: boolean) {
    return patchAndSave({ notify_operations: value })
  }
  function setNotifyServices(value: boolean) {
    return patchAndSave({ notify_services: value })
  }
  function setNotifyErrors(value: boolean) {
    return patchAndSave({ notify_errors: value })
  }
  function setShowToasts(value: boolean) {
    return patchAndSave({ show_toasts: value })
  }
  function setToastPosition(value: ToastPosition) {
    return patchAndSave({ toast_position: value })
  }
  function setKeepNotificationHistory(value: boolean) {
    return patchAndSave({ keep_notification_history: value })
  }
  function setToastDurationMs(value: number) {
    return patchAndSave({ toast_duration_ms: Math.max(1000, value) })
  }
  function setMaxConcurrency(value: number) {
    return patchAndSave({ max_concurrency: Math.max(1, value) })
  }
  function setDebugLog(value: boolean) {
    return patchAndSave({ debug_log: value })
  }
  function setAccentHue(value: number) {
    return patchAndSave({ accent_hue: value })
  }
  function setUiFontSize(value: number) {
    return patchAndSave({ ui_font_size: value })
  }
  function setRowDensity(value: RowDensity) {
    return patchAndSave({ row_density: value })
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
    launchAtLogin,
    showInMenuBar,
    restoreLastPage,
    lastPagePath,
    backupBeforeUpdate,
    cleanupAfterUpdate,
    notifyUpdates,
    notifyOperations,
    notifyServices,
    notifyErrors,
    showToasts,
    toastPosition,
    keepNotificationHistory,
    toastDurationMs,
    maxConcurrency,
    debugLog,
    accentHue,
    uiFontSize,
    rowDensity,
    load,
    setLanguage,
    setTheme,
    setBrewPath,
    setBrewFilePath,
    setAutoUpdate,
    setLogMaxLines,
    setCheckUpdatesOnLaunch,
    setLaunchAtLogin,
    setShowInMenuBar,
    setRestoreLastPage,
    setLastPagePath,
    setBackupBeforeUpdate,
    setCleanupAfterUpdate,
    setNotifyUpdates,
    setNotifyOperations,
    setNotifyServices,
    setNotifyErrors,
    setShowToasts,
    setToastPosition,
    setKeepNotificationHistory,
    setToastDurationMs,
    setMaxConcurrency,
    setDebugLog,
    setAccentHue,
    setUiFontSize,
    setRowDensity,
    resetAll,
    importFromContent,
    validateBrewPath,
  }
})
