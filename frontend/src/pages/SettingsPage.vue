<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Download, RotateCcw, Upload } from 'lucide-vue-next'
import { ConfigPath, Export } from '../../bindings/changeme/services/configservice.js'
import BrewButton from '@/components/common/BrewButton.vue'
import { useLogStore } from '@/stores/log'
import { type AutoUpdateInterval, type ThemeMode, useSettingsStore } from '@/stores/settings'

const { t, locale } = useI18n()
const settingsStore = useSettingsStore()
const logStore = useLogStore()

const draftBrewPath = ref('')
const draftBrewfilePath = ref('')
const importContent = ref('')
const configPath = ref('')
const pathValidationText = ref('')
const pathValidationTone = ref<'muted' | 'success' | 'danger'>('muted')
const busy = ref(false)

const languages = computed(() => [
  { value: 'zh', label: t('settingsCenter.langZh') },
  { value: 'en', label: t('settingsCenter.langEn') },
])
const autoUpdateOptions = computed<Array<{ value: AutoUpdateInterval; label: string }>>(() => [
  { value: 'off', label: t('settingsCenter.autoUpdateOff') },
  { value: '1h', label: '1h' },
  { value: '6h', label: '6h' },
  { value: '12h', label: '12h' },
  { value: '24h', label: '24h' },
])
const themeOptions = computed<Array<{ value: ThemeMode; label: string }>>(() => [
  { value: 'dark', label: t('settings.themeDark') },
  { value: 'light', label: t('settings.themeLight') },
])

const canImport = computed(() => importContent.value.trim().length > 0)

function syncDraft() {
  draftBrewPath.value = settingsStore.brewPath
  draftBrewfilePath.value = settingsStore.brewFilePath
}

async function saveLanguage(value: string) {
  await settingsStore.setLanguage(value)
  locale.value = settingsStore.language
}

async function saveTheme(value: ThemeMode) {
  await settingsStore.setTheme(value)
}

async function saveAutoUpdate(value: AutoUpdateInterval) {
  await settingsStore.setAutoUpdate(value)
}

async function saveCheckUpdates(value: boolean) {
  await settingsStore.setCheckUpdatesOnLaunch(value)
}

async function saveLogMaxLines(value: string) {
  const lines = Math.max(100, Number(value || 1000))
  await settingsStore.setLogMaxLines(lines)
  logStore.setMaxLines(settingsStore.logMaxLines)
}

async function saveBrewPath() {
  await settingsStore.setBrewPath(draftBrewPath.value)
  const result = await settingsStore.validateBrewPath(settingsStore.brewPath)
  if (result.valid) {
    pathValidationTone.value = 'success'
    pathValidationText.value = result.version
      ? `${t('settingsCenter.pathAvailable')} (${result.version})`
      : t('settingsCenter.pathAvailable')
    return
  }
  pathValidationTone.value = 'danger'
  pathValidationText.value = result.error || t('settingsCenter.pathInvalid')
}

async function saveBrewfilePath() {
  await settingsStore.setBrewFilePath(draftBrewfilePath.value)
}

async function exportConfig() {
  const payload = await Export()
  importContent.value = payload
}

async function importConfig() {
  if (!canImport.value) return
  busy.value = true
  try {
    await settingsStore.importFromContent(importContent.value)
    locale.value = settingsStore.language
    logStore.setMaxLines(settingsStore.logMaxLines)
    syncDraft()
  } finally {
    busy.value = false
  }
}

async function resetConfig() {
  busy.value = true
  try {
    await settingsStore.resetAll()
    locale.value = settingsStore.language
    logStore.setMaxLines(settingsStore.logMaxLines)
    syncDraft()
    pathValidationTone.value = 'muted'
    pathValidationText.value = ''
  } finally {
    busy.value = false
  }
}

onMounted(async () => {
  if (!settingsStore.initialized) {
    await settingsStore.load()
  }
  locale.value = settingsStore.language
  logStore.setMaxLines(settingsStore.logMaxLines)
  configPath.value = await ConfigPath()
  syncDraft()
})
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">{{ t('settingsCenter.title') }}</h1>
      <p class="content-subtitle">{{ t('settingsCenter.subtitle') }}</p>
    </div>
    <div class="content-body" style="display:flex; flex-direction:column; gap:16px;">
      <div class="detail-card">
        <div class="detail-card-title">{{ t('settingsCenter.runtime') }}</div>
        <div style="display:grid; gap:12px;">
          <div style="display:flex; justify-content:space-between; align-items:center; gap:16px;">
            <div class="content-subtitle" style="margin:0;">{{ t('settings.language') }}</div>
            <div style="display:flex; gap:8px;">
              <BrewButton v-for="lang in languages" :key="lang.value" :variant="settingsStore.language === lang.value ? 'primary' : 'secondary'" @click="saveLanguage(lang.value)">{{ lang.label }}</BrewButton>
            </div>
          </div>
          <div style="display:flex; justify-content:space-between; align-items:center; gap:16px;">
            <div class="content-subtitle" style="margin:0;">{{ t('settings.theme') }}</div>
            <div style="display:flex; gap:8px;">
              <BrewButton v-for="opt in themeOptions" :key="opt.value" :variant="settingsStore.theme === opt.value ? 'primary' : 'secondary'" @click="saveTheme(opt.value)">{{ opt.label }}</BrewButton>
            </div>
          </div>
          <div>
            <div class="content-subtitle" style="margin-bottom:6px;">{{ t('settingsCenter.brewExecutablePath') }}</div>
            <input v-model="draftBrewPath" class="search-input" :placeholder="t('settingsCenter.brewExecutablePlaceholder')" @change="saveBrewPath">
            <div class="content-subtitle" :style="pathValidationTone === 'success' ? 'color:var(--success);margin-top:6px;' : pathValidationTone === 'danger' ? 'color:var(--danger);margin-top:6px;' : 'margin-top:6px;'">
              {{ pathValidationText || t('settingsCenter.pathApplyHint') }}
            </div>
          </div>
          <div>
            <div class="content-subtitle" style="margin-bottom:6px;">{{ t('settings.brewFilePath') }}</div>
            <input v-model="draftBrewfilePath" class="search-input" :placeholder="t('settings.brewFilePathPlaceholder')" @change="saveBrewfilePath">
          </div>
        </div>
      </div>

      <div class="detail-card">
        <div class="detail-card-title">{{ t('settingsCenter.behavior') }}</div>
        <div style="display:grid; gap:12px;">
          <div style="display:flex; justify-content:space-between; align-items:center; gap:16px;">
            <div class="content-subtitle" style="margin:0;">{{ t('settingsCenter.checkUpdatesOnLaunch') }}</div>
            <label style="display:flex; align-items:center; gap:8px;">
              <input type="checkbox" :checked="settingsStore.checkUpdatesOnLaunch" @change="saveCheckUpdates(($event.target as HTMLInputElement).checked)">
              <span>{{ settingsStore.checkUpdatesOnLaunch ? t('settingsCenter.enabled') : t('settingsCenter.disabled') }}</span>
            </label>
          </div>
          <div style="display:flex; justify-content:space-between; align-items:center; gap:16px;">
            <div class="content-subtitle" style="margin:0;">{{ t('settingsCenter.autoUpdateInterval') }}</div>
            <div style="display:flex; gap:8px; flex-wrap:wrap;">
              <BrewButton v-for="opt in autoUpdateOptions" :key="opt.value" :variant="settingsStore.autoUpdate === opt.value ? 'primary' : 'secondary'" @click="saveAutoUpdate(opt.value)">{{ opt.label }}</BrewButton>
            </div>
          </div>
          <div style="display:flex; justify-content:space-between; align-items:center; gap:16px;">
            <div class="content-subtitle" style="margin:0;">{{ t('settingsCenter.logMaxLines') }}</div>
            <input type="number" class="search-input" style="max-width:180px;" :value="settingsStore.logMaxLines" min="100" step="100" @change="saveLogMaxLines(($event.target as HTMLInputElement).value)">
          </div>
        </div>
      </div>

      <div class="detail-card">
        <div class="detail-card-title">{{ t('settingsCenter.configFile') }}</div>
        <div class="content-subtitle" style="margin-bottom:8px;">{{ configPath }}</div>
        <div class="toolbar" style="padding:0 0 10px; border:0; gap:8px;">
          <BrewButton @click="exportConfig"><Download :size="14" />{{ t('settingsCenter.exportToEditor') }}</BrewButton>
          <BrewButton variant="primary" :disabled="busy || !canImport" @click="importConfig"><Upload :size="14" />{{ t('settingsCenter.importFromEditor') }}</BrewButton>
          <BrewButton variant="danger" :disabled="busy" @click="resetConfig"><RotateCcw :size="14" />{{ t('settingsCenter.resetDefault') }}</BrewButton>
        </div>
        <textarea v-model="importContent" spellcheck="false" style="width:100%; min-height:180px; resize:vertical; background:var(--surface); color:var(--fg); border:1px solid var(--border); border-radius:var(--radius); padding:12px 14px; font-family:var(--font-mono); font-size:12.5px; line-height:1.7; outline:none;" />
      </div>
    </div>
  </section>
</template>
