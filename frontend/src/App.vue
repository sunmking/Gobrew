<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/layout/Sidebar.vue'
import SettingsOverlay from '@/components/layout/SettingsOverlay.vue'
import TitleBar from '@/components/layout/TitleBar.vue'
import NotificationCenter from '@/components/common/NotificationCenter.vue'
import { useInstalledStore } from '@/stores/installed'
import { useLogStore } from '@/stores/log'
import { useSettingsStore } from '@/stores/settings'
import { useUpdateStore } from '@/stores/update'

const { locale } = useI18n()
const settingsStore = useSettingsStore()
const logStore = useLogStore()
const installedStore = useInstalledStore()
const updateStore = useUpdateStore()
const router = useRouter()

onMounted(async () => {
  await settingsStore.load()
  locale.value = settingsStore.language
  logStore.setMaxLines(settingsStore.logMaxLines)
  logStore.setToastDuration(settingsStore.toastDurationMs)
  if (settingsStore.restoreLastPage && settingsStore.lastPagePath && settingsStore.lastPagePath !== router.currentRoute.value.path) {
    await router.replace(settingsStore.lastPagePath)
  }
  await Promise.all([installedStore.fetchInstalled(), updateStore.fetchOutdated()])
})

watch(
  () => router.currentRoute.value.path,
  (path) => {
    if (settingsStore.restoreLastPage) {
      settingsStore.setLastPagePath(path)
    }
  },
)
</script>

<template>
  <div class="window-shell">
    <TitleBar />
    <NotificationCenter />
    <SettingsOverlay />
    <div class="app-body">
      <Sidebar />
      <main class="content-host">
        <router-view v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </router-view>
      </main>
    </div>
  </div>
</template>
