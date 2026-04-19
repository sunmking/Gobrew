<script setup lang="ts">
import { onMounted } from 'vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import { useInstalledStore } from '@/stores/installed'
import { useUpdateStore } from '@/stores/update'

const installedStore = useInstalledStore()
const updateStore = useUpdateStore()

onMounted(async () => {
  await Promise.all([installedStore.fetchInstalled(), updateStore.fetchOutdated()])
})
</script>

<template>
  <div class="app-layout">
    <Sidebar />
    <main class="app-content">
      <router-view v-slot="{ Component }">
        <Transition name="fade" mode="out-in">
          <component :is="Component" />
        </Transition>
      </router-view>
    </main>
  </div>
</template>
