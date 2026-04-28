<script setup lang="ts">
import { Bell, Moon, Settings, Sun } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useLogStore } from '@/stores/log'
import { useSettingsStore } from '@/stores/settings'
import { useUpdateStore } from '@/stores/update'

const router = useRouter()
const { t } = useI18n()
const settingsStore = useSettingsStore()
const updateStore = useUpdateStore()
const logStore = useLogStore()
const isDark = computed(() => settingsStore.theme === 'dark')

async function toggleTheme() {
  await settingsStore.setTheme(isDark.value ? 'light' : 'dark')
}

async function checkUpdates() {
  try {
    await updateStore.forceRefresh()
    const total = updateStore.formulae.length + updateStore.casks.length
    if (total === 0) {
      logStore.markInfo(t('update.upToDate'))
      return
    }
    logStore.markInfo(t('titleBar.updateFound', { count: total }))
  } catch (error: any) {
    logStore.markError(error?.message || t('messages.operationFailed'))
  }
}
</script>

<template>
  <header class="titlebar">
    <div class="titlebar-title">Gobrew</div>
    <div class="titlebar-actions">
      <button class="btn-icon" type="button" :title="isDark ? t('titleBar.switchToLight') : t('titleBar.switchToDark')" @click="toggleTheme">
        <Moon v-if="!isDark" :size="14" />
        <Sun v-else :size="14" />
      </button>
      <button class="btn-icon" type="button" :title="t('titleBar.notifications')" @click="checkUpdates">
        <Bell :size="14" />
      </button>
      <button class="btn-icon" type="button" :title="t('titleBar.settingsCenter')" @click="router.push('/settings')">
        <Settings :size="14" />
      </button>
    </div>
  </header>
</template>
