<script setup lang="ts">
import { Download, Trash2, Upload, MoreHorizontal, Lock, Unlock } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import BrewButton from '@/components/common/BrewButton.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import type { PackageRow } from '@/types/brew'

defineProps<{
  rows: PackageRow[]
  pendingKey?: string
}>()

const emit = defineEmits<{
  (event: 'select', row: PackageRow): void
  (event: 'install', row: PackageRow): void
  (event: 'upgrade', row: PackageRow): void
  (event: 'uninstall', row: PackageRow): void
  (event: 'pin', row: PackageRow): void
  (event: 'unpin', row: PackageRow): void
}>()
const { t } = useI18n()
</script>

<template>
  <table class="pkg-table">
    <thead>
      <tr>
        <th>{{ t('table.type') }}</th>
        <th>{{ t('table.name') }}</th>
        <th>{{ t('table.desc') }}</th>
        <th>{{ t('table.currentVersion') }}</th>
        <th>{{ t('table.latestVersion') }}</th>
        <th>{{ t('table.status') }}</th>
        <th style="width: 150px;"></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in rows" :key="row.key" @click="emit('select', row)">
        <td><span class="pill pill-neutral">{{ row.type === 'cask' ? 'Cask' : 'Formula' }}</span></td>
        <td>
          <div class="pkg-name">{{ row.name }}</div>
          <div class="pkg-version">{{ row.tap }}</div>
        </td>
        <td><div class="pkg-desc">{{ row.desc || t('common.noDetails') }}</div></td>
        <td class="pkg-version">{{ row.installedVersion || '-' }}</td>
        <td class="pkg-version">{{ row.latestVersion || '-' }}</td>
        <td>
          <StatusPill v-if="row.pinned" status="pinned">{{ t('status.pinned') }}</StatusPill>
          <StatusPill v-else-if="row.updateAvailable" status="update">{{ t('status.updatable') }}</StatusPill>
          <StatusPill v-else-if="row.installed" status="installed">{{ t('common.installed') }}</StatusPill>
          <StatusPill v-else status="not-installed">{{ t('status.notInstalled') }}</StatusPill>
        </td>
        <td @click.stop>
          <div style="display:flex; gap:6px; justify-content:flex-end;">
            <BrewButton v-if="!row.installed" variant="primary" :disabled="pendingKey === row.key" @click="emit('install', row)"><Download :size="13" />{{ pendingKey === row.key ? t('actions.installing') : t('actions.install') }}</BrewButton>
            <BrewButton v-else-if="row.updateAvailable && !row.pinned" variant="primary" :disabled="pendingKey === row.key" @click="emit('upgrade', row)"><Upload :size="13" />{{ pendingKey === row.key ? t('actions.updating') : t('actions.update') }}</BrewButton>
            <BrewButton v-else variant="ghost" :disabled="pendingKey === row.key" @click="emit('uninstall', row)"><Trash2 :size="13" />{{ t('actions.uninstall') }}</BrewButton>
            <BrewButton v-if="row.installed && row.type === 'formula' && !row.pinned" variant="ghost" :disabled="pendingKey === row.key" @click="emit('pin', row)"><Lock :size="13" /></BrewButton>
            <BrewButton v-if="row.installed && row.type === 'formula' && row.pinned" variant="ghost" :disabled="pendingKey === row.key" @click="emit('unpin', row)"><Unlock :size="13" /></BrewButton>
            <BrewButton variant="ghost" :disabled="pendingKey === row.key" @click="emit('select', row)"><MoreHorizontal :size="13" /></BrewButton>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>
