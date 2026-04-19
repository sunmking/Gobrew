<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Trash2, ArrowUpCircle, Info, RotateCcw, Download } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useInstalledStore } from '@/stores/installed'
import { useLogStore } from '@/stores/log'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Toast from '@/components/common/Toast.vue'
import BulkActionBar from '@/components/common/BulkActionBar.vue'
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
const installedStore = useInstalledStore()
const logStore = useLogStore()
const selection = useSelection()

const query = ref('')
const action = ref<'uninstall' | 'upgrade'>('uninstall')
const targetKeys = ref<string[]>([])
const confirmOpen = ref(false)
const toastRef = ref<any>(null)
const bulkSummary = ref<BulkSummary | null>(null)
const infoOpen = ref(false)
const infoLoading = ref(false)
const infoError = ref('')
const infoTitle = ref('')
const infoRows = ref<InfoRow[]>([])

function sortByName<T extends { name: string }>(items: T[]) {
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
}

const filteredFormulae = computed(() => {
  const q = query.value.trim().toLowerCase()
  const items = installedStore.formulae.filter((item) => {
    if (!q) return true
    return [item.name, item.full_name, item.desc].join(' ').toLowerCase().includes(q)
  })
  return sortByName(items)
})

const filteredCasks = computed(() => {
  const q = query.value.trim().toLowerCase()
  const items = q
    ? installedStore.casks.filter((item) => [item.name, item.full_name, item.desc].join(' ').toLowerCase().includes(q))
    : installedStore.casks
  return sortByName(items)
})

const allVisibleKeys = computed(() => [
  ...filteredFormulae.value.map((item) => `f:${item.name}`),
  ...filteredCasks.value.map((item) => `c:${item.name}`),
])

const allSelected = computed(() => selection.isAllSelected(allVisibleKeys.value))
const confirmTitle = computed(() => (action.value === 'uninstall' ? t('installed.uninstall') : t('installed.upgrade')))
const confirmMessage = computed(() => {
  if (targetKeys.value.length <= 1) {
    const name = (targetKeys.value[0] || '').split(':')[1] || ''
    return action.value === 'uninstall' ? t('installed.confirmUninstall', { name }) : t('installed.confirmUpgradeSingle', { name })
  }
  return action.value === 'uninstall'
    ? t('installed.confirmUninstallMulti', { count: targetKeys.value.length })
    : t('installed.confirmUpgradeMulti', { count: targetKeys.value.length })
})

watch(
  allVisibleKeys,
  (ids) => {
    selection.sync(ids)
  },
  { immediate: true },
)

function keyToName(key: string) {
  return key.split(':')[1] || key
}

async function refreshInstalledMetadata() {
  await installedStore.forceRefresh()
}

function askAction(nextAction: 'uninstall' | 'upgrade', keys: string[]) {
  const uniqueKeys = [...new Set(keys)]
  if (uniqueKeys.length === 0) return
  action.value = nextAction
  targetKeys.value = uniqueKeys
  confirmOpen.value = true
}

function askUninstall(name: string, isCask = false) {
  askAction('uninstall', [`${isCask ? 'c' : 'f'}:${name}`])
}

function askUpgrade(name: string, isCask = false) {
  askAction('upgrade', [`${isCask ? 'c' : 'f'}:${name}`])
}

async function reinstall(name: string, isCask = false) {
  logStore.startListening()
  try {
    await installedStore.reinstall(name, isCask)
    toastRef.value?.show('success', t('messages.reinstalledName', { name }))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.reinstallFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function install(name: string) {
  logStore.startListening()
  try {
    await BrewService.Install(name)
    toastRef.value?.show('success', t('messages.installedName', { name }))
    await refreshInstalledMetadata()
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.installFailed'))
  } finally {
    logStore.stopListening()
  }
}

function setInfoRows(rows: InfoRow[]) {
  infoRows.value = rows.filter((row) => String(row.value || '').trim() !== '')
}

async function openFormulaInfo(name: string, item: any) {
  infoTitle.value = `${name} · ${t('installed.info')}`
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
    const version = data?.stable_version || data?.installed?.[0]?.version || ''
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
      { label: 'Version', value: item?.versions?.stable || item?.installed?.[0]?.version || '' },
      { label: 'Tap', value: item?.tap || '' },
      { label: 'Description', value: item?.desc || '' },
    ])
    infoError.value = error?.message || ''
  } finally {
    infoLoading.value = false
  }
}

function openCaskInfo(item: any) {
  infoTitle.value = `${item.name} · ${t('installed.info')}`
  infoError.value = ''
  infoLoading.value = false
  infoOpen.value = true
  setInfoRows([
    { label: 'Name', value: item?.name || '' },
    { label: 'Full Name', value: item?.full_name || '' },
    { label: 'Version', value: item?.version || item?.installed || '' },
    { label: 'Tap', value: item?.tap || '' },
    { label: 'Homepage', value: item?.homepage || '' },
    { label: 'Description', value: item?.desc || '' },
  ])
}

async function runAction() {
  confirmOpen.value = false
  const keys = [...new Set(targetKeys.value)]
  if (keys.length === 0) return

  logStore.startListening()
  try {
    const result = await runBulk(keys, async (key) => {
      const [kind, name] = key.split(':')
      const isCask = kind === 'c'
      if (!name) return
      if (action.value === 'uninstall') {
        if (isCask) {
          await BrewService.UninstallCask(name, false)
        } else {
          await BrewService.Uninstall(name)
        }
      } else {
        await BrewService.Upgrade(name)
      }
    })

    bulkSummary.value = {
      action: action.value === 'uninstall' ? t('messages.bulkUninstall') : t('messages.bulkUpgrade'),
      total: result.total,
      success: result.success,
      failures: result.failures.map((item) => ({
        ...item,
        item: keyToName(item.item),
      })),
      timestamp: Date.now(),
    }

    await refreshInstalledMetadata()
    selection.clear()

    if (result.failures.length === 0) {
      toastRef.value?.show('success', t('messages.bulkProcessedPackages', { count: result.success }))
    } else {
      toastRef.value?.show('error', t('messages.bulkFailedPackages', { count: result.failures.length }))
    }
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.operationFailed'))
  } finally {
    targetKeys.value = []
    logStore.stopListening()
  }
}

onMounted(() => {
  installedStore.fetchInstalled()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div v-if="!embedded" class="content-header">
      <div style="display:flex; align-items:center; justify-content:space-between;">
        <h1 class="content-title">{{ t('installed.title') }}</h1>
        <button
          style="background:transparent; border:none; border-radius:var(--radius-sm); padding:6px; cursor:pointer; display:inline-flex; align-items:center; color:var(--color-text-tertiary);"
          @click="refreshInstalledMetadata()"
        >
          <RefreshCw :size="14" />
        </button>
      </div>
      <input
        v-model="query"
        :placeholder="t('installed.search')"
        style="width:100%; margin-top:12px; padding:8px 12px; font-size:13px; border-radius:var(--radius-sm); border:1px solid var(--color-border); background:var(--color-card); color:var(--color-text); outline:none;"
      />
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <LoadingSkeleton v-if="installedStore.loading" :rows="5" />

      <div v-else-if="installedStore.error" style="padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">
        {{ installedStore.error }}
      </div>

      <div v-else-if="filteredFormulae.length === 0 && filteredCasks.length === 0" style="text-align:center; padding:40px 0; font-size:13px; color:var(--color-text-tertiary);">
        {{ t('common.noPackages') }}
      </div>

      <div v-else>
        <BulkActionBar
          :selected-count="selection.selectedCount.value"
          :all-selected="allSelected"
          @select-all="selection.selectAll(allVisibleKeys)"
          @clear="selection.clear()"
        >
          <button
            class="btn-primary"
            :disabled="selection.selectedCount.value === 0"
            @click="askAction('upgrade', selection.selected.value)"
          >
            {{ t('installed.upgrade') }}
          </button>
          <button
            class="btn-danger"
            :disabled="selection.selectedCount.value === 0"
            @click="askAction('uninstall', selection.selected.value)"
          >
            {{ t('installed.uninstall') }}
          </button>
        </BulkActionBar>

        <div v-if="filteredFormulae.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.formulae') }}</h3>
          <div class="card-group">
            <div v-for="item in filteredFormulae" :key="`formula-${item.name}`" class="card-row">
              <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
                <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(`f:${item.name}`)" @change="selection.toggle(`f:${item.name}`)" />
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.versions?.stable || item.installed?.[0]?.version }}</div>
                  <div v-if="item.desc" style="font-size:12px; color:var(--color-text-tertiary); margin-top:2px;">{{ item.desc }}</div>
                </div>
              </div>
              <div style="display:flex; align-items:center; gap:2px;">
                <button class="icon-btn" @click="openFormulaInfo(item.name, item)"><Info :size="15" /></button>
                <button class="icon-btn success" @click="install(item.name)"><Download :size="15" /></button>
                <button class="icon-btn" @click="reinstall(item.name)"><RotateCcw :size="15" /></button>
                <button class="icon-btn accent" @click="askUpgrade(item.name)"><ArrowUpCircle :size="15" /></button>
                <button class="icon-btn danger" @click="askUninstall(item.name)"><Trash2 :size="15" /></button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="filteredCasks.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.casks') }}</h3>
          <div class="card-group">
            <div v-for="item in filteredCasks" :key="`cask-${item.name}`" class="card-row">
              <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
                <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(`c:${item.name}`)" @change="selection.toggle(`c:${item.name}`)" />
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.version || item.installed }}</div>
                  <div v-if="item.desc" style="font-size:12px; color:var(--color-text-tertiary); margin-top:2px;">{{ item.desc }}</div>
                </div>
              </div>
              <div style="display:flex; align-items:center; gap:2px;">
                <button class="icon-btn" @click="openCaskInfo(item)"><Info :size="15" /></button>
                <button class="icon-btn success" @click="install(item.name)"><Download :size="15" /></button>
                <button class="icon-btn" @click="reinstall(item.name, true)"><RotateCcw :size="15" /></button>
                <button class="icon-btn accent" @click="askUpgrade(item.name, true)"><ArrowUpCircle :size="15" /></button>
                <button class="icon-btn danger" @click="askUninstall(item.name, true)"><Trash2 :size="15" /></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :open="confirmOpen"
      :title="confirmTitle"
      :message="confirmMessage"
      :confirm-label="confirmTitle"
      :severity="action === 'uninstall' ? 'danger' : 'default'"
      @confirm="runAction"
      @cancel="confirmOpen = false"
    />
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
.icon-btn.success {
  color: var(--color-success);
}
.icon-btn.success:hover {
  background: var(--color-success-light);
}
.icon-btn.danger {
  color: var(--color-danger);
}
.icon-btn.danger:hover {
  background: var(--color-danger-light);
}
</style>
