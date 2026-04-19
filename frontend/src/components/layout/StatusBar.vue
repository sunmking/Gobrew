<script setup lang="ts">
import { useI18n } from 'vue-i18n'

defineProps<{
  installedCount: number
  outdatedCount: number
  ready: boolean
}>()

const { t } = useI18n()
</script>

<template>
  <footer class="shell-statusbar">
    <div class="shell-status-group">
      <span class="shell-status-item">
        <span aria-hidden="true">📦</span>
        <strong>{{ installedCount }}</strong>
        <span>{{ t('home.installed') }}</span>
      </span>
      <span class="shell-status-dot" :class="{ 'is-syncing': !ready, 'is-warning': ready && outdatedCount > 0 }" />
      <span v-if="!ready" class="shell-status-item">{{ t('statusBar.syncing') }}</span>
      <span v-else-if="outdatedCount > 0" class="shell-status-item">
        <strong>{{ outdatedCount }}</strong>
        <span>{{ t('home.outdated') }}</span>
      </span>
      <span v-else class="shell-status-item">{{ t('update.upToDate') }}</span>
    </div>
    <span class="shell-status-meta">{{ t('statusBar.inventory') }}</span>
  </footer>
</template>
