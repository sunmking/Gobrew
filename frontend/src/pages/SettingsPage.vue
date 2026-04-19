<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Globe, Terminal, Clock, Palette } from 'lucide-vue-next'
import { useSettingsStore, type ThemeMode, type AutoUpdateInterval } from '@/stores/settings'
import Toast from '@/components/common/Toast.vue'

const { t, locale } = useI18n()
const settingsStore = useSettingsStore()
const toastRef = ref<any>(null)

const languages = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
]

const autoUpdateOptions: { value: AutoUpdateInterval; label: string }[] = [
  { value: 'off', label: 'Off' },
  { value: '1h', label: '1h' },
  { value: '6h', label: '6h' },
  { value: '12h', label: '12h' },
  { value: '24h', label: '24h' },
]

const themeOptions: { value: ThemeMode; label: string }[] = [
  { value: 'dark', label: t('settings.themeDark') },
  { value: 'light', label: t('settings.themeLight') },
]

function onLanguageChange() {
  settingsStore.setLanguage(settingsStore.language)
  locale.value = settingsStore.language
}
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header">
      <h1 class="content-title">{{ t('settings.title') }}</h1>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto" style="display:flex;flex-direction:column;gap:20px;">

      <div class="content-section">
        <div class="card-group">
          <div class="card-row" style="padding:14px 16px;">
            <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
              <div style="width:28px;height:28px;border-radius:6px;display:flex;align-items:center;justify-content:center;background:var(--color-accent-light);">
                <Globe :size="15" style="color:var(--color-accent);" />
              </div>
              <div style="min-width:0;">
                <div style="font-size:13px;font-weight:500;">{{ t('settings.language') }}</div>
                <div style="font-size:12px;color:var(--color-text-secondary);margin-top:1px;">{{ t('settings.languageDesc') }}</div>
              </div>
            </div>
            <div style="display:flex;gap:6px;flex-shrink:0;">
              <button
                v-for="lang in languages" :key="lang.value"
                :style="settingsStore.language === lang.value
                  ? 'background:var(--color-accent);color:white;border:none;border-radius:var(--radius-sm);padding:5px 14px;font-size:12px;font-weight:500;cursor:pointer;'
                  : 'background:var(--color-card-hover);color:var(--color-text);border:1px solid var(--color-border);border-radius:var(--radius-sm);padding:5px 14px;font-size:12px;cursor:pointer;'"
                @click="settingsStore.setLanguage(lang.value); onLanguageChange()"
              >{{ lang.label }}</button>
            </div>
          </div>

          <div class="card-row" style="padding:14px 16px;">
            <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
              <div style="width:28px;height:28px;border-radius:6px;display:flex;align-items:center;justify-content:center;background:var(--color-accent-light);">
                <Palette :size="15" style="color:var(--color-accent);" />
              </div>
              <div style="min-width:0;">
                <div style="font-size:13px;font-weight:500;">{{ t('settings.theme') }}</div>
                <div style="font-size:12px;color:var(--color-text-secondary);margin-top:1px;">{{ t('settings.themeDesc') }}</div>
              </div>
            </div>
            <div style="display:flex;gap:6px;flex-shrink:0;">
              <button
                v-for="opt in themeOptions" :key="opt.value"
                :style="settingsStore.theme === opt.value
                  ? 'background:var(--color-accent);color:white;border:none;border-radius:var(--radius-sm);padding:5px 14px;font-size:12px;font-weight:500;cursor:pointer;'
                  : 'background:var(--color-card-hover);color:var(--color-text);border:1px solid var(--color-border);border-radius:var(--radius-sm);padding:5px 14px;font-size:12px;cursor:pointer;'"
                @click="settingsStore.setTheme(opt.value)"
              >{{ opt.label }}</button>
            </div>
          </div>
        </div>
      </div>

      <div class="content-section">
        <div class="section-title">{{ t('settings.autoUpdate') || 'Auto Update' }}</div>
        <div class="card-group">
          <div class="card-row" style="padding:14px 16px;">
            <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
              <div style="width:28px;height:28px;border-radius:6px;display:flex;align-items:center;justify-content:center;background:var(--color-accent-light);">
                <Clock :size="15" style="color:var(--color-accent);" />
              </div>
              <div style="min-width:0;">
                <div style="font-size:13px;font-weight:500;">{{ t('settings.autoUpdate') }}</div>
                <div style="font-size:12px;color:var(--color-text-secondary);margin-top:1px;">{{ t('settings.autoUpdateDesc') }}</div>
              </div>
            </div>
            <div style="display:flex;gap:4px;flex-shrink:0;flex-wrap:wrap;justify-content:flex-end;">
              <button
                v-for="opt in autoUpdateOptions" :key="opt.value"
                :style="settingsStore.autoUpdate === opt.value
                  ? 'background:var(--color-accent);color:white;border:none;border-radius:var(--radius-sm);padding:5px 12px;font-size:12px;font-weight:500;cursor:pointer;'
                  : 'background:var(--color-card-hover);color:var(--color-text);border:1px solid var(--color-border);border-radius:var(--radius-sm);padding:5px 12px;font-size:12px;cursor:pointer;'"
                @click="settingsStore.setAutoUpdate(opt.value)"
              >{{ opt.label }}</button>
            </div>
          </div>
        </div>
      </div>

      <div class="content-section">
        <div class="section-title">{{ t('settings.brewPath') || 'Brew Path' }}</div>
        <div class="card-group">
          <div class="card-row" style="padding:14px 16px;">
            <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
              <div style="width:28px;height:28px;border-radius:6px;display:flex;align-items:center;justify-content:center;background:var(--color-accent-light);">
                <Terminal :size="15" style="color:var(--color-accent);" />
              </div>
              <div style="min-width:0;flex:1;">
                <div style="font-size:13px;font-weight:500;">{{ t('settings.brewPath') }}</div>
                <div style="font-size:12px;color:var(--color-text-secondary);margin-top:1px;margin-bottom:8px;">{{ t('settings.brewPathDesc') }}</div>
                <input
                  v-model="settingsStore.brewPath"
                  type="text"
                  :placeholder="t('settings.brewPathPlaceholder')"
                  style="width:100%;padding:6px 10px;font-size:13px;border-radius:var(--radius-sm);border:1px solid var(--color-border);background:var(--color-content-bg,#f5f5f7);color:var(--color-text);outline:none;font-family:var(--font-mono);"
                  @change="settingsStore.setBrewPath(settingsStore.brewPath)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

    <Toast ref="toastRef" />
  </div>
</template>
