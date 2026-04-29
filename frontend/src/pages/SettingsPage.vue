<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLogStore } from '@/stores/log'
import { type AutoUpdateInterval, type ThemeMode, useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'

type Tab = 'general' | 'appearance' | 'updates' | 'notifications' | 'advanced'
type Prefs = {
  launchAtLogin: boolean
  menuBar: boolean
  restorePage: boolean
  rowDensity: 'compact' | 'default' | 'comfortable'
  backupBeforeUpdate: boolean
  cleanupAfterUpdate: boolean
  notifyUpdates: boolean
  notifyOperations: boolean
  notifyServices: boolean
  notifyErrors: boolean
  showToasts: boolean
  keepNotificationHistory: boolean
  toastPosition: 'tl' | 'tr' | 'bl' | 'br'
  toastDuration: '2000' | '2500' | '4000'
  maxConcurrency: '1' | '4' | '8'
  debugLog: boolean
}

const tabs = [
  { key: 'general' as Tab, icon: '⌘' },
  { key: 'appearance' as Tab, icon: '◐' },
  { key: 'updates' as Tab, icon: '↻' },
  { key: 'notifications' as Tab, icon: '◎' },
  { key: 'advanced' as Tab, icon: '⌥' },
]

const defaultPrefs: Prefs = {
  launchAtLogin: true,
  menuBar: true,
  restorePage: true,
  rowDensity: 'default',
  backupBeforeUpdate: true,
  cleanupAfterUpdate: false,
  notifyUpdates: true,
  notifyOperations: true,
  notifyServices: true,
  notifyErrors: true,
  showToasts: true,
  keepNotificationHistory: true,
  toastPosition: 'br',
  toastDuration: '2500',
  maxConcurrency: '4',
  debugLog: false,
}

const { locale, t } = useI18n()
const settingsStore = useSettingsStore()
const logStore = useLogStore()
const uiStore = useUiStore()

const activeTab = ref<Tab>('general')
const prefs = ref<Prefs>({ ...defaultPrefs })
const lastSavedPrefs = ref<Prefs>({ ...defaultPrefs })
const draftBrewfilePath = ref('')
const draftBrewPath = ref('')
const pathHint = ref(t('settingsPrefs.pathAutoDetectHint'))
const accentHue = ref(255)
const lastSavedAccentHue = ref(255)
const fontSize = ref(13)
const lastSavedFontSize = ref(13)
const themeMode = ref<'light' | 'auto' | 'dark'>('auto')
let saveNoticeTimer: ReturnType<typeof setTimeout> | null = null
let savingPrefs = false
let saveQueued = false

function notifySaveSuccess() {
  if (saveNoticeTimer) clearTimeout(saveNoticeTimer)
  saveNoticeTimer = setTimeout(() => {
    logStore.markSuccess(t('settingsPrefs.messages.saved'))
    saveNoticeTimer = null
  }, 240)
}

function notifySaveError(error: unknown) {
  logStore.markError(
    `${t('settingsPrefs.messages.saveFailed')}: ${error instanceof Error ? error.message : String(error)}`,
  )
}

const autoUpdateFrequency = computed({
  get: () => {
    const v = settingsStore.autoUpdate
    if (v === '24h') return 'daily'
    if (v === 'off') return 'startup'
    return 'weekly'
  },
  set: async (value: 'startup' | 'daily' | 'weekly') => {
    const next: AutoUpdateInterval = value === 'daily' ? '24h' : value === 'weekly' ? '12h' : 'off'
    await settingsStore.setAutoUpdate(next)
    savePrefs()
  },
})

function loadPrefs() {
  prefs.value = {
    launchAtLogin: settingsStore.launchAtLogin,
    menuBar: settingsStore.showInMenuBar,
    restorePage: settingsStore.restoreLastPage,
    rowDensity: settingsStore.rowDensity,
    backupBeforeUpdate: settingsStore.backupBeforeUpdate,
    cleanupAfterUpdate: settingsStore.cleanupAfterUpdate,
    notifyUpdates: settingsStore.notifyUpdates,
    notifyOperations: settingsStore.notifyOperations,
    notifyServices: settingsStore.notifyServices,
    notifyErrors: settingsStore.notifyErrors,
    showToasts: settingsStore.showToasts,
    keepNotificationHistory: settingsStore.keepNotificationHistory,
    toastPosition: settingsStore.toastPosition,
    toastDuration: String(settingsStore.toastDurationMs) as Prefs['toastDuration'],
    maxConcurrency: String(settingsStore.maxConcurrency) as Prefs['maxConcurrency'],
    debugLog: settingsStore.debugLog,
  }
  accentHue.value = settingsStore.accentHue
  fontSize.value = settingsStore.uiFontSize
  lastSavedPrefs.value = { ...prefs.value }
  lastSavedAccentHue.value = accentHue.value
  lastSavedFontSize.value = fontSize.value
}

async function savePrefs() {
  if (savingPrefs) {
    saveQueued = true
    return
  }
  savingPrefs = true

  const tasks: Array<Promise<void>> = []
  const prev = lastSavedPrefs.value
  const next = prefs.value

  if (next.launchAtLogin !== prev.launchAtLogin) tasks.push(settingsStore.setLaunchAtLogin(next.launchAtLogin))
  if (next.menuBar !== prev.menuBar) tasks.push(settingsStore.setShowInMenuBar(next.menuBar))
  if (next.restorePage !== prev.restorePage) tasks.push(settingsStore.setRestoreLastPage(next.restorePage))
  if (next.rowDensity !== prev.rowDensity) tasks.push(settingsStore.setRowDensity(next.rowDensity))
  if (next.backupBeforeUpdate !== prev.backupBeforeUpdate) tasks.push(settingsStore.setBackupBeforeUpdate(next.backupBeforeUpdate))
  if (next.cleanupAfterUpdate !== prev.cleanupAfterUpdate) tasks.push(settingsStore.setCleanupAfterUpdate(next.cleanupAfterUpdate))
  if (next.notifyUpdates !== prev.notifyUpdates) tasks.push(settingsStore.setNotifyUpdates(next.notifyUpdates))
  if (next.notifyOperations !== prev.notifyOperations) tasks.push(settingsStore.setNotifyOperations(next.notifyOperations))
  if (next.notifyServices !== prev.notifyServices) tasks.push(settingsStore.setNotifyServices(next.notifyServices))
  if (next.notifyErrors !== prev.notifyErrors) tasks.push(settingsStore.setNotifyErrors(next.notifyErrors))
  if (next.showToasts !== prev.showToasts) tasks.push(settingsStore.setShowToasts(next.showToasts))
  if (next.keepNotificationHistory !== prev.keepNotificationHistory) tasks.push(settingsStore.setKeepNotificationHistory(next.keepNotificationHistory))
  if (next.toastPosition !== prev.toastPosition) tasks.push(settingsStore.setToastPosition(next.toastPosition))
  if (next.toastDuration !== prev.toastDuration) tasks.push(settingsStore.setToastDurationMs(Number(next.toastDuration)))
  if (next.maxConcurrency !== prev.maxConcurrency) tasks.push(settingsStore.setMaxConcurrency(Number(next.maxConcurrency)))
  if (next.debugLog !== prev.debugLog) tasks.push(settingsStore.setDebugLog(next.debugLog))
  if (accentHue.value !== lastSavedAccentHue.value) tasks.push(settingsStore.setAccentHue(accentHue.value))
  if (fontSize.value !== lastSavedFontSize.value) tasks.push(settingsStore.setUiFontSize(fontSize.value))

  try {
    if (tasks.length > 0) {
      await Promise.all(tasks)
      lastSavedPrefs.value = { ...prefs.value }
      lastSavedAccentHue.value = accentHue.value
      lastSavedFontSize.value = fontSize.value
      logStore.setToastDuration(settingsStore.toastDurationMs)
      notifySaveSuccess()
    }
  } catch (error) {
    notifySaveError(error)
  } finally {
    savingPrefs = false
    if (saveQueued) {
      saveQueued = false
      void savePrefs()
    }
  }
}

function applyAppearance() {
  document.documentElement.style.setProperty('--accent', `oklch(58% 0.18 ${accentHue.value})`)
  document.documentElement.style.setProperty('--accent-bg', `oklch(95% 0.03 ${accentHue.value})`)
  document.documentElement.style.setProperty('--ui-font-size', `${fontSize.value}px`)
  document.documentElement.style.setProperty('--row-density', prefs.value.rowDensity)
}

async function setTheme(mode: 'light' | 'auto' | 'dark') {
  themeMode.value = mode
  if (mode === 'auto') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    await settingsStore.setTheme(prefersDark ? 'dark' : 'light')
  } else {
    await settingsStore.setTheme(mode as ThemeMode)
  }
  savePrefs()
}

async function setLanguage(value: 'zh' | 'en') {
  try {
    await settingsStore.setLanguage(value)
    locale.value = settingsStore.language
    notifySaveSuccess()
  } catch (error) {
    notifySaveError(error)
  }
}

async function saveBrewfilePath() {
  try {
    await settingsStore.setBrewFilePath(draftBrewfilePath.value)
    notifySaveSuccess()
  } catch (error) {
    notifySaveError(error)
  }
}

async function saveBrewPath() {
  try {
    await settingsStore.setBrewPath(draftBrewPath.value)
    const result = await settingsStore.validateBrewPath(settingsStore.brewPath)
    pathHint.value = result.valid
      ? t('settingsPrefs.pathAvailable', { version: result.version ? ` (${result.version})` : '' })
      : (result.error || t('settingsPrefs.pathInvalid'))
    notifySaveSuccess()
  } catch (error) {
    notifySaveError(error)
  }
}

async function resetAll() {
  await settingsStore.resetAll()
  loadPrefs()
  draftBrewfilePath.value = settingsStore.brewFilePath
  draftBrewPath.value = settingsStore.brewPath
  themeMode.value = 'auto'
  applyAppearance()
  locale.value = settingsStore.language
  logStore.setMaxLines(settingsStore.logMaxLines)
  logStore.setToastDuration(settingsStore.toastDurationMs)
  pathHint.value = t('settingsPrefs.pathAutoDetectHint')
  logStore.markInfo(t('settingsPrefs.messages.reset'))
}

onMounted(async () => {
  if (!settingsStore.initialized) await settingsStore.load()
  loadPrefs()
  locale.value = settingsStore.language
  draftBrewfilePath.value = settingsStore.brewFilePath
  draftBrewPath.value = settingsStore.brewPath
  themeMode.value = settingsStore.theme
  applyAppearance()
})
</script>

<template>
  <aside class="prefs-sidebar">
          <div class="prefs-title">{{ t('settingsPrefs.title') }}</div>
          <button v-for="tab in tabs" :key="tab.key" class="prefs-nav" :class="{ active: activeTab === tab.key }" @click="activeTab = tab.key">
            <span class="prefs-icon">{{ tab.icon }}</span>
            <span>{{ t(`settingsPrefs.tabs.${tab.key}`) }}</span>
          </button>
  </aside>

  <main class="prefs-content">
          <section v-show="activeTab === 'general'" class="prefs-page active">
            <div class="prefs-section">{{ t('settingsPrefs.sections.startup') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.startup.launchAtLogin') }}</div></div><label class="toggle"><input v-model="prefs.launchAtLogin" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.startup.menuBar') }}</div></div><label class="toggle"><input v-model="prefs.menuBar" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.startup.restorePage') }}</div><div class="prefs-help">{{ t('settingsPrefs.startup.restorePageHelp') }}</div></div><label class="toggle"><input v-model="prefs.restorePage" type="checkbox" @change="savePrefs"><span /></label></div>
            </div>
            <div class="prefs-section">{{ t('settingsPrefs.sections.languageAndPath') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.language') }}</div></div><select class="prefs-select" :value="settingsStore.language" @change="setLanguage(($event.target as HTMLSelectElement).value as 'zh' | 'en')"><option value="zh">{{ t('settingsPrefs.langZh') }}</option><option value="en">{{ t('settingsPrefs.langEn') }}</option></select></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.brewfilePath') }}</div><div class="prefs-help">{{ t('settingsPrefs.brewfilePathHelp') }}</div></div><input v-model="draftBrewfilePath" class="prefs-input" @change="saveBrewfilePath"></div>
            </div>
          </section>

          <section v-show="activeTab === 'appearance'" class="prefs-page active">
            <div class="prefs-section">{{ t('settingsPrefs.sections.theme') }}</div>
            <div class="prefs-group">
              <div class="prefs-row">
                <div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.appearance.colorScheme') }}</div></div>
                <div class="segmented">
                  <button :class="{ active: themeMode === 'light' }" @click="setTheme('light')">{{ t('settingsPrefs.appearance.themeLight') }}</button>
                  <button :class="{ active: themeMode === 'auto' }" @click="setTheme('auto')">{{ t('settingsPrefs.appearance.themeAuto') }}</button>
                  <button :class="{ active: themeMode === 'dark' }" @click="setTheme('dark')">{{ t('settingsPrefs.appearance.themeDark') }}</button>
                </div>
              </div>
              <div class="prefs-row">
                <div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.appearance.accentColor') }}</div></div>
                <div class="accent-swatches">
                  <button class="accent-swatch" :class="{ active: accentHue === 255 }" style="background:oklch(58% 0.18 255)" @click="accentHue = 255; applyAppearance(); savePrefs()" />
                  <button class="accent-swatch" :class="{ active: accentHue === 145 }" style="background:oklch(58% 0.16 145)" @click="accentHue = 145; applyAppearance(); savePrefs()" />
                  <button class="accent-swatch" :class="{ active: accentHue === 35 }" style="background:oklch(58% 0.18 35)" @click="accentHue = 35; applyAppearance(); savePrefs()" />
                  <button class="accent-swatch" :class="{ active: accentHue === 310 }" style="background:oklch(58% 0.18 310)" @click="accentHue = 310; applyAppearance(); savePrefs()" />
                </div>
              </div>
            </div>
            <div class="prefs-section">{{ t('settingsPrefs.sections.displayDensity') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.appearance.fontSize') }}</div><div class="prefs-help">{{ fontSize }} px{{ fontSize === 13 ? t('settingsPrefs.appearance.defaultSizeSuffix') : '' }}</div></div><input v-model.number="fontSize" class="prefs-range" type="range" min="12" max="16" @input="applyAppearance" @change="savePrefs"></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.appearance.rowDensity') }}</div></div><select v-model="prefs.rowDensity" class="prefs-select" @change="applyAppearance(); savePrefs()"><option value="compact">{{ t('settingsPrefs.appearance.rowDensityCompact') }}</option><option value="default">{{ t('settingsPrefs.appearance.rowDensityDefault') }}</option><option value="comfortable">{{ t('settingsPrefs.appearance.rowDensityComfortable') }}</option></select></div>
            </div>
          </section>

          <section v-show="activeTab === 'updates'" class="prefs-page active">
            <div class="prefs-section">{{ t('settingsPrefs.sections.autoCheck') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.updates.checkHomebrewUpdates') }}</div></div><label class="toggle"><input :checked="settingsStore.checkUpdatesOnLaunch" type="checkbox" @change="settingsStore.setCheckUpdatesOnLaunch(($event.target as HTMLInputElement).checked)"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.updates.frequency') }}</div></div><select class="prefs-select" :value="autoUpdateFrequency" @change="autoUpdateFrequency = ($event.target as HTMLSelectElement).value as 'startup' | 'daily' | 'weekly'"><option value="startup">{{ t('settingsPrefs.updates.frequencyStartup') }}</option><option value="daily">{{ t('settingsPrefs.updates.frequencyDaily') }}</option><option value="weekly">{{ t('settingsPrefs.updates.frequencyWeekly') }}</option></select></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.updates.backupBeforeUpdate') }}</div><div class="prefs-help">{{ t('settingsPrefs.updates.backupBeforeUpdateHelp') }}</div></div><label class="toggle"><input v-model="prefs.backupBeforeUpdate" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.updates.cleanupAfterUpdate') }}</div></div><label class="toggle"><input v-model="prefs.cleanupAfterUpdate" type="checkbox" @change="savePrefs"><span /></label></div>
            </div>
          </section>

          <section v-show="activeTab === 'notifications'" class="prefs-page active">
            <div class="prefs-section">{{ t('settingsPrefs.sections.notifications') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.packageUpdates') }}</div></div><label class="toggle"><input v-model="prefs.notifyUpdates" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.operationsDone') }}</div></div><label class="toggle"><input v-model="prefs.notifyOperations" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.serviceChanges') }}</div></div><label class="toggle"><input v-model="prefs.notifyServices" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.errorChanges') }}</div></div><label class="toggle"><input v-model="prefs.notifyErrors" type="checkbox" @change="savePrefs"><span /></label></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.toastPosition') }}</div></div><select v-model="prefs.toastPosition" class="prefs-select" @change="savePrefs"><option value="br">{{ t('settingsPrefs.notifications.positionBr') }}</option><option value="bl">{{ t('settingsPrefs.notifications.positionBl') }}</option><option value="tr">{{ t('settingsPrefs.notifications.positionTr') }}</option><option value="tl">{{ t('settingsPrefs.notifications.positionTl') }}</option></select></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.toastDuration') }}</div></div><select v-model="prefs.toastDuration" class="prefs-select" @change="savePrefs"><option value="2000">{{ t('settingsPrefs.notifications.toast2s') }}</option><option value="2500">{{ t('settingsPrefs.notifications.toast25s') }}</option><option value="4000">{{ t('settingsPrefs.notifications.toast4s') }}</option></select></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.showToasts') }}</div></div><label class="toggle"><input v-model="prefs.showToasts" type="checkbox" @change="savePrefs"><span /></label></div>
            </div>
            <div class="prefs-section">{{ t('settingsPrefs.notifications.centerSection') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.notifications.keepHistory') }}</div></div><label class="toggle"><input v-model="prefs.keepNotificationHistory" type="checkbox" @change="savePrefs"><span /></label></div>
            </div>
          </section>

          <section v-show="activeTab === 'advanced'" class="prefs-page active">
            <div class="prefs-section">{{ t('settingsPrefs.sections.homebrew') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.advanced.brewBinary') }}</div><div class="prefs-help">{{ pathHint }}</div></div><input v-model="draftBrewPath" class="prefs-input" :placeholder="t('settingsPrefs.advanced.brewBinaryPlaceholder')" @change="saveBrewPath"></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.advanced.maxConcurrency') }}</div></div><select v-model="prefs.maxConcurrency" class="prefs-select" @change="savePrefs"><option value="1">1</option><option value="4">4</option><option value="8">8</option></select></div>
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.advanced.debugLog') }}</div><div class="prefs-help">{{ t('settingsPrefs.advanced.debugLogHelp') }}</div></div><label class="toggle"><input v-model="prefs.debugLog" type="checkbox" @change="savePrefs"><span /></label></div>
            </div>
            <div class="prefs-section">{{ t('settingsPrefs.sections.dangerZone') }}</div>
            <div class="prefs-group">
              <div class="prefs-row"><div class="prefs-row-main"><div class="prefs-label">{{ t('settingsPrefs.advanced.resetAll') }}</div><div class="prefs-help">{{ t('settingsPrefs.advanced.resetAllHelp') }}</div></div><button class="btn btn-danger" @click="resetAll">{{ t('settingsPrefs.advanced.reset') }}</button></div>
            </div>
          </section>
  </main>
  <footer class="prefs-footer">
    <button class="btn btn-secondary" @click="uiStore.closeSettings()">{{ t('common.cancel') }}</button>
    <button class="btn btn-primary" @click="uiStore.closeSettings(); logStore.markSuccess(t('settingsPrefs.messages.applied'))">{{ t('settingsPrefs.done') }}</button>
  </footer>
</template>

<style scoped>
.prefs-sidebar { width: 156px; flex-shrink: 0; background: var(--surface); border-right: 1px solid var(--border); padding: 14px 8px; }
.prefs-title { padding: 0 8px 10px; border-bottom: 1px solid var(--border); margin-bottom: 8px; font-size: 12px; font-weight: 600; letter-spacing: 0.01em; }
.prefs-nav { width: 100%; border: 0; background: transparent; display: flex; align-items: center; gap: 8px; padding: 7px 8px; border-radius: var(--radius); color: var(--muted); cursor: pointer; text-align: left; font-size: 12px; font-weight: 500; }
.prefs-nav:hover { background: var(--border); color: var(--fg); }
.prefs-nav.active { background: var(--accent-bg); color: var(--accent); font-weight: 500; }
.prefs-nav:focus-visible { outline: 2px solid var(--accent); outline-offset: 1px; }
.prefs-icon { width: 18px; text-align: center; font-size: 12px; }
.prefs-content { flex: 1; overflow: auto; padding: 18px 20px 72px; }
.prefs-page { display: block; }
.prefs-section { margin: 18px 0 8px; font-size: 11px; font-weight: 700; letter-spacing: 0.07em; text-transform: uppercase; color: var(--muted); }
.prefs-section:first-child { margin-top: 0; }
.prefs-group { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); overflow: hidden; }
.prefs-row { min-height: 46px; display: flex; align-items: center; gap: 12px; padding: 10px 14px; border-bottom: 1px solid var(--border); }
.prefs-row:last-child { border-bottom: 0; }
.prefs-row-main { flex: 1; min-width: 0; }
.prefs-label { font-size: 12.5px; font-weight: 500; color: var(--fg); line-height: 1.35; }
.prefs-help { margin-top: 2px; font-size: 11.5px; color: var(--muted); }
.prefs-input, .prefs-select { height: 30px; min-width: 160px; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--surface); color: var(--fg); font-size: 12.5px; padding: 0 9px; outline: none; }
.prefs-input { min-width: 210px; font-family: var(--font-mono); }
.prefs-input:focus, .prefs-select:focus { border-color: var(--accent); box-shadow: 0 0 0 2px oklch(58% 0.18 255 / 0.12); }
.toggle { position: relative; width: 36px; height: 22px; flex-shrink: 0; }
.toggle input { position: absolute; opacity: 0; inset: 0; }
.toggle span { display: block; width: 36px; height: 22px; border-radius: 999px; background: var(--border); cursor: pointer; transition: background 0.16s ease; }
.toggle span::after { content: ""; position: absolute; top: 3px; left: 3px; width: 16px; height: 16px; border-radius: 50%; background: white; box-shadow: 0 1px 3px oklch(0% 0 0 / 0.25); transition: transform 0.16s ease; }
.toggle input:checked + span { background: var(--accent); }
.toggle input:checked + span::after { transform: translateX(14px); }
.toggle input:focus-visible + span { outline: 2px solid var(--accent); outline-offset: 2px; }
.segmented { display: flex; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius-sm); }
.segmented button { height: 30px; padding: 0 11px; border: 0; border-right: 1px solid var(--border); background: transparent; color: var(--muted); cursor: pointer; font-size: 12px; font-weight: 500; }
.segmented button:last-child { border-right: 0; }
.segmented button.active { background: var(--accent); color: white; }
.segmented button:focus-visible { outline: 2px solid var(--accent); outline-offset: -2px; position: relative; z-index: 1; }
.accent-swatches { display: flex; gap: 8px; }
.accent-swatch { width: 20px; height: 20px; border: 0; border-radius: 50%; cursor: pointer; box-shadow: 0 0 0 1px oklch(0% 0 0 / 0.12); }
.accent-swatch.active { outline: 2px solid var(--accent); outline-offset: 2px; }
.accent-swatch:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; }
.prefs-range { width: 128px; accent-color: var(--accent); }
.prefs-range:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; border-radius: 4px; }
.prefs-footer { position: absolute; left: 156px; right: 0; bottom: 0; height: 52px; padding: 10px 20px; display: flex; align-items: center; justify-content: flex-end; gap: 8px; background: var(--surface); border-top: 1px solid var(--border); }
</style>
