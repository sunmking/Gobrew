import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'dark' | 'light'
export type AutoUpdateInterval = 'off' | '1h' | '6h' | '12h' | '24h'

function loadSetting(key: string, fallback: string): string {
  return localStorage.getItem(key) || fallback
}

function applyTheme(t: ThemeMode) {
  const root = document.documentElement
  if (!root) return
  root.classList.toggle('dark', t === 'dark')
}

export const useSettingsStore = defineStore('settings', () => {
  const language = ref(loadSetting('gobrew-lang', navigator.language.startsWith('zh') ? 'zh' : 'en'))
  const brewPath = ref(loadSetting('gobrew-brew-path', ''))
  const brewFilePath = ref(loadSetting('gobrew-brewfile-path', ''))
  const autoUpdate = ref<AutoUpdateInterval>(loadSetting('gobrew-auto-update', 'off') as AutoUpdateInterval)
  const theme = ref<ThemeMode>(loadSetting('gobrew-theme', 'dark') as ThemeMode)

  watch(theme, (t) => applyTheme(t), { immediate: true })

  function setLanguage(lang: string) {
    language.value = lang
    localStorage.setItem('gobrew-lang', lang)
  }

  function setBrewPath(path: string) {
    brewPath.value = path
    localStorage.setItem('gobrew-brew-path', path)
  }

  function setBrewFilePath(path: string) {
    brewFilePath.value = path
    localStorage.setItem('gobrew-brewfile-path', path)
  }

  function setAutoUpdate(interval: AutoUpdateInterval) {
    autoUpdate.value = interval
    localStorage.setItem('gobrew-auto-update', interval)
  }

  function setTheme(t: ThemeMode) {
    theme.value = t
    localStorage.setItem('gobrew-theme', t)
  }

  return {
    language, brewPath, brewFilePath, autoUpdate, theme,
    setLanguage, setBrewPath, setBrewFilePath, setAutoUpdate, setTheme,
  }
})
