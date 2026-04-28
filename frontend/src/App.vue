<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Sidebar from '@/components/layout/Sidebar.vue'
import TitleBar from '@/components/layout/TitleBar.vue'
import OperationStatusBar from '@/components/common/OperationStatusBar.vue'
import { useInstalledStore } from '@/stores/installed'
import { useLogStore } from '@/stores/log'
import { useSettingsStore } from '@/stores/settings'
import { useUpdateStore } from '@/stores/update'

const { locale } = useI18n()
const settingsStore = useSettingsStore()
const logStore = useLogStore()
const installedStore = useInstalledStore()
const updateStore = useUpdateStore()

onMounted(async () => {
  await settingsStore.load()
  locale.value = settingsStore.language
  logStore.setMaxLines(settingsStore.logMaxLines)
  await Promise.all([installedStore.fetchInstalled(), updateStore.fetchOutdated()])
})
</script>

<template>
  <div class="window-shell">
    <TitleBar />
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
    <OperationStatusBar />
  </div>
</template>
