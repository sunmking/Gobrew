<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Search, Download, Info, Check } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useSearchStore } from '@/stores/search'
import { useLogStore } from '@/stores/log'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import Toast from '@/components/common/Toast.vue'
import BulkResultSummary from '@/components/common/BulkResultSummary.vue'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'
import PackageInfoDialog from '@/components/common/PackageInfoDialog.vue'
import PackageIcon from '@/components/icons/PackageIcon.vue'
import type { InfoRow } from '@/components/common/PackageInfoDialog.vue'

const props = withDefaults(
  defineProps<{
    embedded?: boolean
  }>(),
  {
    embedded: false,
  },
)

const { t } = useI18n()
const searchStore = useSearchStore()
const logStore = useLogStore()
const selection = useSelection()
const toastRef = ref<any>(null)
const query = ref('')
const bulkSummary = ref<BulkSummary | null>(null)
const infoOpen = ref(false)
const infoLoading = ref(false)
const infoError = ref('')
const infoTitle = ref('')
const infoRows = ref<InfoRow[]>([])

const installing = ref<Set<string>>(new Set())
const DEFAULT_RENDER_COUNT = 80
const RENDER_STEP = 80
const formulaVisibleCount = ref(DEFAULT_RENDER_COUNT)
const caskVisibleCount = ref(DEFAULT_RENDER_COUNT)

function sortByName<T extends { name: string }>(items: T[]) {
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
}

const sortedFormulae = computed(() => sortByName(searchStore.results.formulae))
const sortedCasks = computed(() => sortByName(searchStore.results.casks))
const formulae = computed(() => sortedFormulae.value.slice(0, formulaVisibleCount.value))
const casks = computed(() => sortedCasks.value.slice(0, caskVisibleCount.value))
const hasMoreFormulae = computed(() => formulaVisibleCount.value < sortedFormulae.value.length)
const hasMoreCasks = computed(() => caskVisibleCount.value < sortedCasks.value.length)
const allVisibleKeys = computed(() => [...formulae.value.map((item) => `f:${item.name}`), ...casks.value.map((item) => `c:${item.name}`)])
const allSelected = computed(() => selection.isAllSelected(allVisibleKeys.value))
const selectedNames = computed(() => [...new Set(selection.selected.value.map((key) => key.split(':')[1]).filter(Boolean))])

watch(
  allVisibleKeys,
  (ids) => {
    selection.sync(ids)
  },
  { immediate: true },
)

watch(
  () => [searchStore.query, sortedFormulae.value.length, sortedCasks.value.length],
  () => {
    formulaVisibleCount.value = DEFAULT_RENDER_COUNT
    caskVisibleCount.value = DEFAULT_RENDER_COUNT
  },
)

async function doSearch() {
  selection.clear()
  await searchStore.search(query.value)
}

async function install(name: string) {
  installing.value.add(name)
  logStore.startListening()
  try {
    await searchStore.install(name)
    toastRef.value?.show('success', `${name} ${t('explore.installed')}`)
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || `${t('explore.install')} failed`)
  } finally {
    logStore.stopListening()
    installing.value.delete(name)
  }
}

async function installSelected() {
  const names = selectedNames.value
  if (names.length === 0) return

  logStore.startListening()
  try {
    const result = await runBulk(names, (name) => BrewService.Install(name))
    bulkSummary.value = {
      action: 'Bulk Install',
      total: result.total,
      success: result.success,
      failures: result.failures,
      timestamp: Date.now(),
    }
    selection.clear()
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', `${result.success}/${result.total} installed`)
  } finally {
    logStore.stopListening()
  }
}

function selectAllVisible() {
  selection.selectAll(allVisibleKeys.value)
}

function loadMoreFormulae() {
  formulaVisibleCount.value += RENDER_STEP
}

function loadMoreCasks() {
  caskVisibleCount.value += RENDER_STEP
}

function setInfoRows(rows: InfoRow[]) {
  infoRows.value = rows.filter((row) => String(row.value || '').trim() !== '')
}

async function openFormulaInfo(name: string, item: any) {
  infoTitle.value = `${name} · ${t('explore.info')}`
  infoError.value = ''
  infoLoading.value = true
  infoOpen.value = true
  try {
    const [info, deps, uses] = await Promise.all([
      BrewService.Info(name),
      BrewService.Deps(name, true),
      BrewService.Uses(name, true, true),
    ])
    const data: any = info
    const version = data?.stable_version || data?.versions?.stable || data?.installed?.[0]?.version || ''
    setInfoRows([
      { label: 'Name', value: data?.name || name },
      { label: 'Full Name', value: data?.full_name || '' },
      { label: 'Version', value: version },
      { label: 'Tap', value: data?.tap || '' },
      { label: 'License', value: data?.license || '' },
      { label: 'Homepage', value: data?.homepage || '' },
      { label: 'Description', value: data?.desc || '' },
      { label: 'Dependencies (tree)', value: (deps || []).join('\n') },
      { label: 'Used By (installed)', value: (uses || []).join('\n') },
    ])
  } catch (error: any) {
    setInfoRows([
      { label: 'Name', value: item?.name || name },
      { label: 'Full Name', value: item?.full_name || '' },
      { label: 'Tap', value: item?.tap || '' },
      { label: 'Description', value: item?.desc || '' },
    ])
    infoError.value = error?.message || ''
  } finally {
    infoLoading.value = false
  }
}

function openCaskInfo(item: any) {
  infoTitle.value = `${item.name} · ${t('explore.info')}`
  infoError.value = ''
  infoLoading.value = false
  infoOpen.value = true
  setInfoRows([
    { label: 'Name', value: item?.name || '' },
    { label: 'Full Name', value: item?.full_name || '' },
    { label: 'Tap', value: item?.tap || '' },
    { label: 'Description', value: item?.desc || '' },
  ])
}

const hasSelection = computed(() => selection.selectedCount.value > 0)

onMounted(() => {
  // Keep search results ephemeral: always start from a clean state when entering the page.
  searchStore.clearResults()
  query.value = ''
})

onUnmounted(() => {
  searchStore.clearResults()
  query.value = ''
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div v-if="!embedded" class="content-header">
      <h1 class="content-title">{{ t('explore.title') }}</h1>
    </div>
    <form class="ui-toolbar" style="flex-shrink:0;justify-content:space-between;" @submit.prevent="doSearch">
      <div class="ui-toolbar" style="gap:6px;">
        <input
          class="ui-input"
          v-model="query"
          :placeholder="t('explore.search')"
          style="width:min(36vw,260px);min-width:170px;height:30px;padding:0 10px;"
        />
        <button
          type="submit"
          class="ui-btn ui-btn-primary ui-btn-sm"
          style="height:30px;"
        >
          <Search :size="12" />
          {{ t('explore.searchBtn') }}
        </button>
      </div>
      <div class="ui-toolbar" style="gap:6px;margin-left:auto;">
        <button
          type="button"
          class="ui-btn ui-btn-secondary ui-btn-sm"
          style="height:30px;"
          :disabled="allVisibleKeys.length === 0 || allSelected"
          @click="selectAllVisible"
        >
          {{ t('bulk.selectAll') }}
        </button>
        <button
          type="button"
          class="ui-btn ui-btn-secondary ui-btn-sm"
          style="height:30px;"
          :disabled="!hasSelection"
          @click="selection.clear()"
        >
          {{ t('bulk.clear') }}
        </button>
        <button
          type="button"
          class="ui-btn ui-btn-primary ui-btn-sm"
          style="height:30px;"
          :disabled="!hasSelection"
          @click="installSelected"
        >
          {{ t('explore.install') }} ({{ selection.selectedCount.value }})
        </button>
      </div>
    </form>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <LoadingSkeleton v-if="searchStore.loading" :rows="5" />

      <div v-else-if="searchStore.error" style="padding:12px;border-radius:var(--radius-sm);background:var(--color-danger-light);color:var(--color-danger);font-size:13px;">
        {{ searchStore.error }}
      </div>

      <div v-else-if="formulae.length === 0 && casks.length === 0" style="text-align:center;padding:40px 0;font-size:13px;color:var(--color-text-tertiary);">
        {{ t('explore.noResults') }}
      </div>

      <div v-else>
        <div v-if="formulae.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.formulae') }} ({{ formulae.length }}/{{ sortedFormulae.length }})</h3>
          <div class="card-group">
            <div v-for="item in formulae" :key="`formula-${item.name}`" class="card-row" :style="selection.isSelected(`f:${item.name}`) ? 'background:var(--color-accent-light);' : ''">
              <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
                <div
                  style="width:18px;height:18px;border-radius:4px;border:2px solid var(--color-border);display:flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:all 100ms ease;"
                  :style="selection.isSelected(`f:${item.name}`) ? 'background:var(--color-accent);border-color:var(--color-accent);' : ''"
                  @click="selection.toggle(`f:${item.name}`)"
                >
                  <Check v-if="selection.isSelected(`f:${item.name}`)" :size="12" style="color:white;" />
                </div>
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.full_name }}</div>
                  <div v-if="item.desc" style="font-size:12px;color:var(--color-text-tertiary);margin-top:2px;">{{ item.desc }}</div>
                </div>
              </div>
              <div style="display:flex;align-items:center;gap:4px;">
                <button class="icon-btn" @click="openFormulaInfo(item.name, item)">
                  <Info :size="16" />
                </button>
                <button class="icon-btn accent" :disabled="installing.has(item.name)" @click="install(item.name)">
                  <Download :size="16" />
                </button>
              </div>
            </div>
          </div>
          <div v-if="hasMoreFormulae" style="display:flex;justify-content:center;margin-top:8px;">
            <button class="ui-btn ui-btn-secondary ui-btn-sm" @click="loadMoreFormulae">
              {{ t('common.loadMore') || 'Load more' }}
            </button>
          </div>
        </div>

        <div v-if="casks.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.casks') }} ({{ casks.length }}/{{ sortedCasks.length }})</h3>
          <div class="card-group">
            <div v-for="item in casks" :key="`cask-${item.name}`" class="card-row" :style="selection.isSelected(`c:${item.name}`) ? 'background:var(--color-accent-light);' : ''">
              <div style="display:flex;align-items:center;gap:10px;min-width:0;flex:1;">
                <div
                  style="width:18px;height:18px;border-radius:4px;border:2px solid var(--color-border);display:flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:all 100ms ease;"
                  :style="selection.isSelected(`c:${item.name}`) ? 'background:var(--color-accent);border-color:var(--color-accent);' : ''"
                  @click="selection.toggle(`c:${item.name}`)"
                >
                  <Check v-if="selection.isSelected(`c:${item.name}`)" :size="12" style="color:white;" />
                </div>
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.full_name }}</div>
                  <div v-if="item.desc" style="font-size:12px;color:var(--color-text-tertiary);margin-top:2px;">{{ item.desc }}</div>
                </div>
              </div>
              <div style="display:flex;align-items:center;gap:4px;">
                <button class="icon-btn" @click="openCaskInfo(item)">
                  <Info :size="16" />
                </button>
                <button class="icon-btn accent" :disabled="installing.has(item.name)" @click="install(item.name)">
                  <Download :size="16" />
                </button>
              </div>
            </div>
          </div>
          <div v-if="hasMoreCasks" style="display:flex;justify-content:center;margin-top:8px;">
            <button class="ui-btn ui-btn-secondary ui-btn-sm" @click="loadMoreCasks">
              {{ t('common.loadMore') || 'Load more' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <PackageInfoDialog
      :open="infoOpen"
      :title="infoTitle"
      :rows="infoRows"
      :loading="infoLoading"
      :error="infoError"
      @close="infoOpen = false"
    />
    <Toast ref="toastRef" />
  </div>
</template>

<style scoped>
.icon-btn {
  background: transparent;
  border: none;
  padding: 4px;
  cursor: pointer;
  border-radius: 4px;
  color: var(--color-text-tertiary);
  display: inline-flex;
  align-items: center;
}
.icon-btn:hover {
  background: var(--color-accent-light);
  color: var(--color-accent);
}
.icon-btn.accent {
  color: var(--color-accent);
}
.icon-btn.accent:hover {
  background: var(--color-accent-light);
}
.icon-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}
</style>
