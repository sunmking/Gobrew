<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
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
const { t } = useI18n()
const installedStore = useInstalledStore()
const updateStore = useUpdateStore()
const searchStore = useSearchStore()
const logStore = useLogStore()
const filter = ref('all')
const query = ref('')
const pendingKey = ref('')
const pageSize = 50
const visibleLimit = ref(pageSize)
const scrollRoot = ref<HTMLElement | null>(null)
let searchTimer: ReturnType<typeof setTimeout> | null = null

const filterOptions = computed(() => [
  { label: t('allPackages.filterAll'), value: 'all' },
  { label: t('allPackages.filterInstalled'), value: 'installed' },
  { label: t('allPackages.filterUpdates'), value: 'updates' },
])

function formulaVersion(item: any) {
  return item.installed?.[0]?.version || item.stable_version || item.versions?.stable || ''
}

function upsert(rows: Map<string, PackageRow>, row: PackageRow) {
  const current = rows.get(row.key)
  rows.set(row.key, current ? { ...current, ...row, installed: current.installed || row.installed, updateAvailable: current.updateAvailable || row.updateAvailable, pinned: current.pinned || row.pinned } : row)
}

const rows = computed(() => {
  const map = new Map<string, PackageRow>()
  const hasQuery = query.value.trim().length > 0
  const installedFormulae = new Map(installedStore.formulae.map(item => [item.name, item]))
  const installedCasks = new Map(installedStore.casks.map(item => [item.token || item.name, item]))
  const outdatedFormulae = new Map(updateStore.formulae.map(item => [item.name, item]))
  const outdatedCasks = new Map(updateStore.casks.map(item => [item.name, item]))

  if (!hasQuery) for (const item of installedStore.formulae) {
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
      pinned: item.pinned,
    })
  }
  if (!hasQuery) for (const item of installedStore.casks) {
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
      pinned: false,
    })
  }
  if (!hasQuery) for (const item of updateStore.formulae) {
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
      pinned: item.pinned,
    })
  }
  if (!hasQuery) for (const item of updateStore.casks) {
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
      pinned: false,
    })
  }
  for (const item of searchStore.results.formulae) {
    const installed = installedFormulae.get(item.name)
    const outdated = outdatedFormulae.get(item.name)
    upsert(map, {
      key: `formula:${item.name}`,
      type: 'formula',
      name: item.name,
      fullName: item.full_name || item.name,
      desc: item.desc || installed?.desc || '',
      tap: item.tap || installed?.tap || '',
      installedVersion: installed ? formulaVersion(installed) : outdated?.installed_versions?.join(', ') || '',
      latestVersion: outdated?.current_version || installed?.stable_version || installed?.versions?.stable || '',
      installed: Boolean(installed || outdated),
      updateAvailable: Boolean(outdated),
      pinned: Boolean(installed?.pinned || outdated?.pinned),
    })
  }
  for (const item of searchStore.results.casks) {
    const installed = installedCasks.get(item.name)
    const outdated = outdatedCasks.get(item.name)
    upsert(map, {
      key: `cask:${item.name}`,
      type: 'cask',
      name: item.name,
      fullName: item.full_name || item.name,
      desc: item.desc || installed?.desc || '',
      tap: item.tap || installed?.tap || '',
      installedVersion: installed?.installed || installed?.version || outdated?.installed_version || '',
      latestVersion: outdated?.current_version || installed?.version || '',
      installed: Boolean(installed || outdated),
      updateAvailable: Boolean(outdated),
      pinned: false,
    })
  }

  return Array.from(map.values())
    .filter(row => filter.value === 'installed' ? row.installed : filter.value === 'updates' ? row.updateAvailable : true)
    .sort((a, b) => Number(b.updateAvailable) - Number(a.updateAvailable) || a.name.localeCompare(b.name))
})

const displayRows = computed(() => rows.value.slice(0, visibleLimit.value))
const hasMoreRows = computed(() => visibleLimit.value < rows.value.length)

watch(query, value => {
  if (searchTimer) clearTimeout(searchTimer)
  if (value.trim()) filter.value = 'all'
  searchTimer = setTimeout(() => searchStore.search(value), 260)
})

watch([rows, filter, query], () => {
  visibleLimit.value = pageSize
  nextTick(loadMoreIfNeeded)
})

function loadMoreRows() {
  if (!hasMoreRows.value) return
  visibleLimit.value = Math.min(visibleLimit.value + pageSize, rows.value.length)
  nextTick(loadMoreIfNeeded)
}

function loadMoreIfNeeded() {
  const el = scrollRoot.value
  if (!el || !hasMoreRows.value) return
  if (el.scrollHeight <= el.clientHeight + 160) {
    loadMoreRows()
  }
}

function onScroll() {
  const el = scrollRoot.value
  if (!el || !hasMoreRows.value) return
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 180) {
    loadMoreRows()
  }
}

async function refresh() {
  await Promise.all([installedStore.forceRefresh(), updateStore.forceRefresh()])
  if (query.value.trim()) await searchStore.search(query.value)
}

function select(row: PackageRow) {
  router.push(`/packages/${row.type}/${encodeURIComponent(row.name)}`)
}

async function runAction(action: 'install' | 'upgrade' | 'uninstall' | 'pin' | 'unpin', row: PackageRow) {
  if (pendingKey.value) return
  pendingKey.value = row.key
  logStore.startListening(`${action} ${row.name}`)
  try {
    if (action === 'install') await BrewService.InstallPackage(row.name, row.type)
    if (action === 'upgrade') await BrewService.Upgrade(row.name)
    if (action === 'pin') await BrewService.Pin(row.name)
    if (action === 'unpin') await BrewService.Unpin(row.name)
    if (action === 'uninstall') {
      if (row.type === 'cask') await BrewService.UninstallCask(row.name, false)
      else await BrewService.Uninstall(row.name)
    }
    await refresh()
  } finally {
    pendingKey.value = ''
    logStore.stopListening()
  }
}

onMounted(async () => {
  await refresh()
  nextTick(loadMoreIfNeeded)
})

onUnmounted(() => {
  if (searchTimer) clearTimeout(searchTimer)
})
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <ToolbarSearch v-model="query" :placeholder="t('allPackages.searchPlaceholder')" @submit="searchStore.search(query)" />
      <div class="toolbar-spacer">
        <SegmentedControl v-model="filter" :options="filterOptions" />
        <BrewButton @click="refresh"><RefreshCw :size="14" />{{ t('common.refresh') }}</BrewButton>
      </div>
    </div>
    <div v-if="installedStore.error || updateStore.error || searchStore.error" class="content-body">
      <p class="error-text">{{ installedStore.error || updateStore.error || searchStore.error }}</p>
    </div>
    <div v-else ref="scrollRoot" class="content-body" style="padding-top:0;" @scroll="onScroll">
      <div v-if="searchStore.loading" class="content-subtitle" style="padding:10px 12px;">{{ t('allPackages.searching') }}</div>
      <PackageTable
        :rows="displayRows"
        :pending-key="pendingKey"
        @select="select"
        @install="runAction('install', $event)"
        @upgrade="runAction('upgrade', $event)"
        @uninstall="runAction('uninstall', $event)"
        @pin="runAction('pin', $event)"
        @unpin="runAction('unpin', $event)"
      />
      <div v-if="rows.length === 0 && !searchStore.loading" class="empty-state">{{ query.trim() ? t('allPackages.noMatch') : t('allPackages.noRows') }}</div>
      <div v-else-if="hasMoreRows" class="empty-state" style="padding:14px;">{{ t('allPackages.scrollMore', { shown: displayRows.length, total: rows.length }) }}</div>
      <div v-else-if="rows.length > pageSize" class="empty-state" style="padding:14px;">{{ t('allPackages.showAll', { total: rows.length }) }}</div>
    </div>
  </section>
</template>
