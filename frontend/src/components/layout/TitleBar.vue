<script setup lang="ts">
import { Bell, Moon, Settings, Sun } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '@/stores/notification'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'

const { t } = useI18n()
const settingsStore = useSettingsStore()
const uiStore = useUiStore()
const notificationStore = useNotificationStore()
const isDark = computed(() => settingsStore.theme === 'dark')

async function toggleTheme() {
  await settingsStore.setTheme(isDark.value ? 'light' : 'dark')
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
      <button class="btn-icon notif-button" type="button" :title="t('titleBar.notifications')" @click.stop="notificationStore.toggle()">
        <Bell :size="14" />
        <span class="notif-dot" :class="{ hidden: notificationStore.unreadCount === 0 }" />
      </button>
      <button class="btn-icon" type="button" :title="t('titleBar.settingsCenter')" @click="uiStore.openSettings()">
        <Settings :size="14" />
      </button>
    </div>
  </header>
</template>

<style scoped>
.notif-button { position: relative; }
.notif-dot {
  position:absolute; top:5px; right:5px;
  width:8px; height:8px;
  border-radius:999px; background:var(--danger);
  box-shadow:0 0 0 2px var(--surface);
}
.notif-dot.hidden { display:none; }
</style>
