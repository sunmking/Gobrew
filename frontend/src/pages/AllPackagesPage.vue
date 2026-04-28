<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { RefreshCw } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import ToolbarSearch from '@/components/common/ToolbarSearch.vue'
import SegmentedControl from '@/components/common/SegmentedControl.vue'
import BrewButton from '@/components/common/BrewButton.vue'
import PackageTable from '@/components/packages/PackageTable.vue'
import { useInstalledStore } from '@/stores/installed'
import { useUpdateStore } from '@/stores/update'
import { useSearchStore } from '@/stores/search'
import { useLogStore } from '@/stores/log'
import type { PackageRow, PackageType } from '@/types/brew'

const router = useRouter()
const installedStore = useInstalledStore()
const updateStore = useUpdateStore()
const searchStore = useSearchStore()
const logStore = useLogStore()
const filter = ref('all')
const query = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null

const filterOptions = [
  { label: '全部', value: 'all' },
  { label: '已安装', value: 'installed' },
  { label: '可更新', value: 'updates' },
]

function formulaVersion(item: any) {
  return item.installed?.[0]?.version || item.stable_version || item.versions?.stable || ''
}

function upsert(rows: Map<string, PackageRow>, row: PackageRow) {
  const current = rows.get(row.key)
  rows.set(row.key, current ? { ...current, ...row, installed: current.installed || row.installed, updateAvailable: current.updateAvailable || row.updateAvailable } : row)
}

const rows = computed(() => {
  const map = new Map<string, PackageRow>()
  for (const item of installedStore.formulae) {
    const name = item.name
    upsert(map, {
      key: `formula:${name}`,
      type: 'formula',
      name,
      fullName: item.full_name || name,
      desc: item.desc,
      tap: item.tap,
      installedVersion: formulaVersion(item),
      latestVersion: item.stable_version || item.versions?.stable || formulaVersion(item),
      installed: true,
      updateAvailable: false,
    })
  }
  for (const item of installedStore.casks) {
    const name = item.token || item.name
    upsert(map, {
      key: `cask:${name}`,
      type: 'cask',
      name,
      fullName: item.full_name || name,
      desc: item.desc,
      tap: item.tap,
      installedVersion: item.installed || item.version,
      latestVersion: item.version,
      installed: true,
      updateAvailable: false,
    })
  }
  for (const item of updateStore.formulae) {
    upsert(map, {
      key: `formula:${item.name}`,
      type: 'formula',
      name: item.name,
      fullName: item.name,
      desc: '',
      tap: '',
      installedVersion: item.installed_versions?.join(', ') || '',
      latestVersion: item.current_version,
      installed: true,
      updateAvailable: true,
    })
  }
  for (const item of updateStore.casks) {
    upsert(map, {
      key: `cask:${item.name}`,
      type: 'cask',
      name: item.name,
      fullName: item.name,
      desc: '',
      tap: '',
      installedVersion: item.installed_version,
      latestVersion: item.current_version,
      installed: true,
      updateAvailable: true,
    })
  }
  for (const item of searchStore.results.formulae) {
    upsert(map, {
      key: `formula:${item.name}`,
      type: 'formula',
      name: item.name,
      fullName: item.full_name || item.name,
      desc: item.desc,
      tap: item.tap,
      installedVersion: '',
      latestVersion: '',
      installed: false,
      updateAvailable: false,
    })
  }
  for (const item of searchStore.results.casks) {
    upsert(map, {
      key: `cask:${item.name}`,
      type: 'cask',
      name: item.name,
      fullName: item.full_name || item.name,
      desc: item.desc,
      tap: item.tap,
      installedVersion: '',
      latestVersion: '',
      installed: false,
      updateAvailable: false,
    })
  }

  return Array.from(map.values())
    .filter(row => filter.value === 'installed' ? row.installed : filter.value === 'updates' ? row.updateAvailable : true)
    .sort((a, b) => Number(b.updateAvailable) - Number(a.updateAvailable) || a.name.localeCompare(b.name))
})

watch(query, value => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => searchStore.search(value), 260)
})

async function refresh() {
  await Promise.all([installedStore.forceRefresh(), updateStore.forceRefresh()])
  if (query.value.trim()) await searchStore.search(query.value)
}

function select(row: PackageRow) {
  router.push(`/packages/${row.type}/${encodeURIComponent(row.name)}`)
}

async function runAction(action: 'install' | 'upgrade' | 'uninstall', row: PackageRow) {
  logStore.startListening(`${action} ${row.name}`)
  try {
    if (action === 'install') await BrewService.Install(row.name)
    if (action === 'upgrade') await BrewService.Upgrade(row.name)
    if (action === 'uninstall') {
      if (row.type === 'cask') await BrewService.UninstallCask(row.name, false)
      else await BrewService.Uninstall(row.name)
    }
    await refresh()
  } finally {
    logStore.stopListening()
  }
}

onMounted(refresh)
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <ToolbarSearch v-model="query" placeholder="搜索包名称或描述..." @submit="searchStore.search(query)" />
      <div class="toolbar-spacer">
        <SegmentedControl v-model="filter" :options="filterOptions" />
        <BrewButton @click="refresh"><RefreshCw :size="14" />刷新</BrewButton>
      </div>
    </div>
    <div v-if="installedStore.error || updateStore.error || searchStore.error" class="content-body">
      <p class="error-text">{{ installedStore.error || updateStore.error || searchStore.error }}</p>
    </div>
    <div v-else class="content-body" style="padding-top:0;">
      <PackageTable :rows="rows" @select="select" @install="runAction('install', $event)" @upgrade="runAction('upgrade', $event)" @uninstall="runAction('uninstall', $event)" />
      <div v-if="rows.length === 0" class="empty-state">没有匹配的包</div>
    </div>
  </section>
</template>
