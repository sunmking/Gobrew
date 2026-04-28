<script setup lang="ts">
import { Download, Trash2, Upload, MoreHorizontal } from 'lucide-vue-next'
import BrewButton from '@/components/common/BrewButton.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import type { PackageRow } from '@/types/brew'

defineProps<{ rows: PackageRow[] }>()

const emit = defineEmits<{
  (event: 'select', row: PackageRow): void
  (event: 'install', row: PackageRow): void
  (event: 'upgrade', row: PackageRow): void
  (event: 'uninstall', row: PackageRow): void
}>()
</script>

<template>
  <table class="pkg-table">
    <thead>
      <tr>
        <th>类型</th>
        <th>包名</th>
        <th>描述</th>
        <th>当前版本</th>
        <th>最新版本</th>
        <th>状态</th>
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
        <td><div class="pkg-desc">{{ row.desc || '暂无描述' }}</div></td>
        <td class="pkg-version">{{ row.installedVersion || '-' }}</td>
        <td class="pkg-version">{{ row.latestVersion || '-' }}</td>
        <td>
          <StatusPill v-if="row.updateAvailable" status="update">可更新</StatusPill>
          <StatusPill v-else-if="row.installed" status="installed">已安装</StatusPill>
          <StatusPill v-else status="not-installed">未安装</StatusPill>
        </td>
        <td @click.stop>
          <div style="display:flex; gap:6px; justify-content:flex-end;">
            <BrewButton v-if="!row.installed" variant="primary" @click="emit('install', row)"><Download :size="13" />安装</BrewButton>
            <BrewButton v-else-if="row.updateAvailable" variant="primary" @click="emit('upgrade', row)"><Upload :size="13" />更新</BrewButton>
            <BrewButton v-else variant="ghost" @click="emit('uninstall', row)"><Trash2 :size="13" />卸载</BrewButton>
            <BrewButton variant="ghost" @click="emit('select', row)"><MoreHorizontal :size="13" /></BrewButton>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>
