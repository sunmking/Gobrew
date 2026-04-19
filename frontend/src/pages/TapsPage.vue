<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Trash2, RefreshCw, Info, Download, ArrowUpCircle } from 'lucide-vue-next'
import * as TapService from '../../bindings/changeme/services/tapservice.js'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useTapsStore } from '@/stores/taps'
import { useLogStore } from '@/stores/log'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Toast from '@/components/common/Toast.vue'
import BulkActionBar from '@/components/common/BulkActionBar.vue'
import BulkResultSummary from '@/components/common/BulkResultSummary.vue'
import PackageIcon from '@/components/icons/PackageIcon.vue'
import type { TapDetail } from '@/types/brew'

const props = withDefaults(
  defineProps<{
    embedded?: boolean
  }>(),
  {
    embedded: false,
  },
)

const { t } = useI18n()
const tapsStore = useTapsStore()
const logStore = useLogStore()
const selection = useSelection()

const newTapName = ref('')
const removeTarget = ref('')
const showRemoveDialog = ref(false)
const toastRef = ref<any>(null)
const query = ref('')
const bulkSummary = ref<BulkSummary | null>(null)

const infoOpen = ref(false)
const infoTitle = ref('')
const infoLoading = ref(false)
const infoError = ref('')
const infoTapName = ref('')

type TapPackage = {
  id: string
  name: string
  type: 'formula' | 'cask'
  installed: boolean
  outdated: boolean
}

const tapRemote = ref('')
const tapBranch = ref('')
const tapLastCommit = ref('')
const tapCustomRemote = ref(false)
const tapPackages = ref<TapPackage[]>([])
const packageQuery = ref('')
const selectedPackageIds = ref<string[]>([])
const refreshing = ref(false)
const addingTap = ref(false)
const removingTap = ref(false)
const bulkRemoving = ref(false)
const packageActionLoading = ref(false)

const taps = computed(() => {
  const q = query.value.trim().toLowerCase()
  const items = q ? tapsStore.taps.filter((tap) => [tap.name, tap.remote].join(' ').toLowerCase().includes(q)) : tapsStore.taps
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
})

const removableTapNames = computed(() => taps.value.filter((tap) => tap.custom_remote).map((tap) => tap.name))
const allSelected = computed(() => selection.isAllSelected(removableTapNames.value))

const filteredTapPackages = computed(() => {
  const q = packageQuery.value.trim().toLowerCase()
  const items = q ? tapPackages.value.filter((pkg) => pkg.name.toLowerCase().includes(q)) : tapPackages.value
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
})

const formulaPackages = computed(() => filteredTapPackages.value.filter((pkg) => pkg.type === 'formula'))
const caskPackages = computed(() => filteredTapPackages.value.filter((pkg) => pkg.type === 'cask'))
const selectedCount = computed(() => selectedPackageIds.value.length)

watch(
  removableTapNames,
  (ids) => {
    selection.sync(ids)
  },
  { immediate: true },
)

watch(filteredTapPackages, (items) => {
  const ids = new Set(items.map((item) => item.id))
  selectedPackageIds.value = selectedPackageIds.value.filter((id) => ids.has(id))
})

function normalizePackageName(name: string) {
  const trimmed = (name || '').trim()
  if (!trimmed) return ''
  const parts = trimmed.split('/').filter(Boolean)
  return parts[parts.length - 1] || trimmed
}

function collectCandidateNames(raw: string) {
  const full = (raw || '').trim()
  const short = normalizePackageName(full)
  const set = new Set<string>()
  if (full) set.add(full)
  if (short) set.add(short)
  return set
}

function isSelectedPkg(id: string) {
  return selectedPackageIds.value.includes(id)
}

function togglePkg(id: string) {
  if (isSelectedPkg(id)) {
    selectedPackageIds.value = selectedPackageIds.value.filter((x) => x !== id)
    return
  }
  selectedPackageIds.value = [...selectedPackageIds.value, id]
}

function selectAllPkg() {
  selectedPackageIds.value = filteredTapPackages.value.map((item) => item.id)
}

function clearPkgSelection() {
  selectedPackageIds.value = []
}

async function addTap() {
  if (!newTapName.value.trim()) return
  addingTap.value = true
  logStore.startListening()
  try {
    await tapsStore.add(newTapName.value.trim())
    toastRef.value?.show('success', t('messages.tapAdded', { name: newTapName.value }))
    newTapName.value = ''
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.failedAddTap'))
  } finally {
    logStore.stopListening()
    addingTap.value = false
  }
}

async function confirmRemove() {
  showRemoveDialog.value = false
  removingTap.value = true
  logStore.startListening()
  try {
    await tapsStore.remove(removeTarget.value)
    toastRef.value?.show('success', t('messages.tapRemoved', { name: removeTarget.value }))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.failedRemoveTap'))
  } finally {
    logStore.stopListening()
    removingTap.value = false
  }
}

async function removeSelected() {
  const names = [...selection.selected.value]
  if (names.length === 0) return

  bulkRemoving.value = true
  logStore.startListening()
  try {
    const result = await runBulk(names, (name) => TapService.Remove(name))
    bulkSummary.value = {
      action: t('messages.bulkRemoveTaps'),
      total: result.total,
      success: result.success,
      failures: result.failures,
      timestamp: Date.now(),
    }
    await tapsStore.fetchTaps()
    selection.clear()
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', t('messages.bulkRemovedSummary', { success: result.success, total: result.total }))
  } finally {
    logStore.stopListening()
    bulkRemoving.value = false
  }
}

async function refreshTaps() {
  refreshing.value = true
  try {
    await tapsStore.fetchTaps()
  } finally {
    refreshing.value = false
  }
}

async function loadTapDetail(tapName: string, fallback: any) {
  const detail = ((await TapService.Details(tapName)) as unknown) as TapDetail
  tapRemote.value = detail?.remote || fallback?.remote || ''
  tapBranch.value = detail?.branch || ''
  tapLastCommit.value = detail?.last_commit || ''
  tapCustomRemote.value = !!detail?.custom_remote

  const [installedResult, outdatedResult] = await Promise.all([
    BrewService.ListInstalled(),
    BrewService.Outdated(),
  ])

  const installedFormulaNames = new Set<string>()
  const installedCaskNames = new Set<string>()
  for (const item of installedResult?.formulae || []) {
    for (const name of collectCandidateNames(item?.full_name || item?.name || '')) installedFormulaNames.add(name)
  }
  for (const item of installedResult?.casks || []) {
    for (const name of collectCandidateNames(item?.full_name || item?.name || '')) installedCaskNames.add(name)
  }

  const outdatedFormulaNames = new Set<string>()
  const outdatedCaskNames = new Set<string>()
  for (const item of outdatedResult?.formulae || []) {
    for (const name of collectCandidateNames(item?.name || '')) outdatedFormulaNames.add(name)
  }
  for (const item of outdatedResult?.casks || []) {
    for (const name of collectCandidateNames(item?.name || '')) outdatedCaskNames.add(name)
  }

  const formulaItems: TapPackage[] = (detail?.formula_names || []).map((raw) => {
    const candidates = collectCandidateNames(raw)
    const installed = [...candidates].some((name) => installedFormulaNames.has(name))
    const outdated = [...candidates].some((name) => outdatedFormulaNames.has(name))
    return {
      id: `formula:${raw}`,
      name: raw,
      type: 'formula',
      installed,
      outdated,
    }
  })

  const caskItems: TapPackage[] = (detail?.cask_tokens || []).map((raw) => {
    const candidates = collectCandidateNames(raw)
    const installed = [...candidates].some((name) => installedCaskNames.has(name))
    const outdated = [...candidates].some((name) => outdatedCaskNames.has(name))
    return {
      id: `cask:${raw}`,
      name: raw,
      type: 'cask',
      installed,
      outdated,
    }
  })

  tapPackages.value = [...formulaItems, ...caskItems]
}

async function openTapInfo(tap: any) {
  infoTitle.value = `${tap.name} · ${t('explore.info')}`
  infoTapName.value = tap.name
  infoLoading.value = true
  infoError.value = ''
  infoOpen.value = true
  packageQuery.value = ''
  selectedPackageIds.value = []
  tapPackages.value = []

  try {
    await loadTapDetail(tap.name, tap)
  } catch (error: any) {
    infoError.value = error?.message || t('messages.failedLoadTapDetails')
  } finally {
    infoLoading.value = false
  }
}

async function runPackageAction(pkg: TapPackage, action: 'install' | 'uninstall' | 'upgrade') {
  packageActionLoading.value = true
  logStore.startListening()
  try {
    await TapService.PackageAction(pkg.type, pkg.name, action)
    await loadTapDetail(infoTapName.value, { remote: tapRemote.value, custom_remote: tapCustomRemote.value })
    toastRef.value?.show('success', t('messages.packageActionDone', { action, name: pkg.name }))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.packageActionFailed', { action }))
  } finally {
    logStore.stopListening()
    packageActionLoading.value = false
  }
}

async function runBulkPackageAction(action: 'install' | 'uninstall' | 'upgrade') {
  const targets = tapPackages.value.filter((pkg) => selectedPackageIds.value.includes(pkg.id))
  if (targets.length === 0) return

  packageActionLoading.value = true
  logStore.startListening()
  try {
    const result = await runBulk(
      targets.map((pkg) => pkg.id),
      async (id) => {
        const pkg = tapPackages.value.find((item) => item.id === id)
        if (!pkg) return
        await TapService.PackageAction(pkg.type, pkg.name, action)
      },
    )
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', t('messages.bulkActionSummary', { success: result.success, total: result.total, action }))
    await loadTapDetail(infoTapName.value, { remote: tapRemote.value, custom_remote: tapCustomRemote.value })
    selectedPackageIds.value = []
  } finally {
    logStore.stopListening()
    packageActionLoading.value = false
  }
}

onMounted(() => {
  refreshTaps()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header" :style="embedded ? 'margin-bottom:10px;' : ''">
      <div style="display:flex; align-items:center; justify-content:space-between;">
        <h1 v-if="!embedded" class="content-title">{{ t('taps.title') }}</h1>
        <div v-else />
        <button
          class="ui-btn ui-btn-secondary ui-btn-sm"
          :disabled="refreshing || tapsStore.loading"
          @click="refreshTaps()"
        >
          <RefreshCw :size="14" :class="refreshing ? 'is-spinning' : ''" />
          {{ t('common.refresh') }}
        </button>
      </div>
      <form style="display:flex; align-items:center; gap:8px; margin-top:12px;" @submit.prevent="addTap">
        <input
          v-model="newTapName"
          type="text"
          :placeholder="t('taps.repoPlaceholder')"
          class="ui-input"
          style="flex:1;"
        />
        <button
          type="submit"
          class="ui-btn ui-btn-primary ui-btn-sm"
          :disabled="addingTap || !newTapName.trim()"
        >
          <Plus :size="14" :class="addingTap ? 'is-spinning' : ''" />
          {{ addingTap ? `${t('common.loading')}` : t('taps.add') }}
        </button>
      </form>
      <input
        v-model="query"
        type="text"
        :placeholder="t('placeholders.searchTaps')"
        class="ui-input"
        style="width:100%; margin-top:8px;"
      />
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <div v-if="tapsStore.loading" style="font-size:13px; color:var(--color-text-tertiary);">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="tapsStore.error" style="padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">
        {{ tapsStore.error }}
      </div>

      <div v-else-if="taps.length === 0" style="text-align:center; padding:40px 0; font-size:13px; color:var(--color-text-tertiary);">
        {{ t('taps.noTaps') }}
      </div>

      <div v-else>
        <BulkActionBar
          :selected-count="selection.selectedCount.value"
          :all-selected="allSelected"
          @select-all="selection.selectAll(removableTapNames)"
          @clear="selection.clear()"
        >
          <button
            class="btn-danger"
            :disabled="selection.selectedCount.value === 0 || bulkRemoving"
            @click="removeSelected"
          >
            {{ bulkRemoving ? t('common.loading') : t('taps.remove') }}
          </button>
        </BulkActionBar>

        <div class="card-group">
          <div v-for="tap in taps" :key="tap.name" class="card-row">
            <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
              <input
                v-if="tap.custom_remote"
                type="checkbox"
                style="flex-shrink:0;"
                :checked="selection.isSelected(tap.name)"
                @change="selection.toggle(tap.name)"
              />
              <div v-else style="width:16px; flex-shrink:0;" />
              <div style="min-width:0;">
                <div class="card-row-title">{{ tap.name }}</div>
                <div class="card-row-subtitle">{{ tap.remote }}</div>
              </div>
            </div>
            <div style="display:flex; align-items:center; gap:4px;">
              <button class="icon-btn" @click="openTapInfo(tap)">
                <Info :size="15" />
              </button>
              <button
                v-if="tap.custom_remote"
                class="icon-btn danger"
                @click="removeTarget = tap.name; showRemoveDialog = true"
              >
                <Trash2 :size="15" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :open="showRemoveDialog"
      :title="t('taps.remove')"
      :message="t('taps.confirmRemove', { name: removeTarget })"
      :confirm-label="t('taps.remove')"
      severity="danger"
      @confirm="confirmRemove"
      @cancel="showRemoveDialog = false"
    />

    <teleport to="body">
      <div v-if="infoOpen" class="modal-overlay" @click.self="infoOpen = false">
        <div class="modal-content">
          <div style="display:flex; align-items:center; justify-content:space-between; gap:8px;">
            <h3 style="font-size:15px; font-weight:600; color:var(--color-text); margin:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;">{{ infoTitle }}</h3>
            <button
              style="background:transparent; color:var(--color-text-secondary); border:1px solid var(--color-border); border-radius:var(--radius-sm); padding:4px 10px; font-size:12px; cursor:pointer;"
              @click="infoOpen = false"
            >
              {{ t('common.close') }}
            </button>
          </div>

          <div class="stat-grid" style="margin-top:12px; grid-template-columns: repeat(2, 1fr);">
            <div class="stat-card" style="padding:10px 14px;">
              <div style="font-size:11px; color:var(--color-text-tertiary);">{{ t('taps.detail.remote') }}</div>
              <div style="font-size:13px; color:var(--color-text); margin-top:2px;">{{ tapRemote || '-' }}</div>
            </div>
            <div class="stat-card" style="padding:10px 14px;">
              <div style="font-size:11px; color:var(--color-text-tertiary);">{{ t('taps.detail.branch') }}</div>
              <div style="font-size:13px; color:var(--color-text); margin-top:2px;">{{ tapBranch || '-' }}</div>
            </div>
            <div class="stat-card" style="padding:10px 14px;">
              <div style="font-size:11px; color:var(--color-text-tertiary);">{{ t('taps.detail.lastCommit') }}</div>
              <div style="font-size:13px; color:var(--color-text); margin-top:2px;">{{ tapLastCommit || '-' }}</div>
            </div>
            <div class="stat-card" style="padding:10px 14px;">
              <div style="font-size:11px; color:var(--color-text-tertiary);">{{ t('taps.detail.customRemote') }}</div>
              <div style="font-size:13px; color:var(--color-text); margin-top:2px;">{{ tapCustomRemote ? t('common.yes') : t('common.no') }}</div>
            </div>
          </div>

          <div style="margin-top:12px; border:1px solid var(--color-border); border-radius:var(--radius-lg); padding:12px; display:flex; flex-direction:column; flex:1; min-height:0;">
            <div style="display:flex; flex-wrap:wrap; align-items:center; justify-content:space-between; gap:8px;">
              <input
                v-model="packageQuery"
                type="text"
                :placeholder="t('placeholders.filterPackages')"
                style="padding:6px 10px; font-size:13px; border-radius:var(--radius-sm); border:1px solid var(--color-border); background:var(--color-card); color:var(--color-text); outline:none; width:100%; max-width:280px;"
              />
              <div style="display:flex; align-items:center; gap:6px;">
                <button
                  class="ui-btn ui-btn-secondary ui-btn-sm"
                  :disabled="packageActionLoading"
                  @click="selectAllPkg"
                >{{ t('bulk.selectAll') }}</button>
                <button
                  class="ui-btn ui-btn-secondary ui-btn-sm"
                  :disabled="packageActionLoading"
                  @click="clearPkgSelection"
                >{{ t('bulk.clear') }}</button>
                <button
                  class="ui-btn ui-btn-primary ui-btn-sm"
                  :disabled="selectedCount === 0 || packageActionLoading"
                  @click="runBulkPackageAction('install')"
                >{{ t('explore.install') }}</button>
                <button
                  class="ui-btn ui-btn-primary ui-btn-sm"
                  :disabled="selectedCount === 0 || packageActionLoading"
                  @click="runBulkPackageAction('upgrade')"
                >{{ t('installed.upgrade') }}</button>
                <button
                  class="ui-btn ui-btn-danger ui-btn-sm"
                  :disabled="selectedCount === 0 || packageActionLoading"
                  @click="runBulkPackageAction('uninstall')"
                >{{ t('installed.uninstall') }}</button>
              </div>
            </div>

            <div v-if="infoLoading" style="margin-top:12px; font-size:13px; color:var(--color-text-tertiary);">{{ t('common.loading') }}</div>
            <div v-else-if="infoError" style="margin-top:12px; padding:10px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">{{ infoError }}</div>
            <div v-else-if="filteredTapPackages.length === 0" style="margin-top:12px; font-size:13px; color:var(--color-text-tertiary);">{{ t('common.noPackages') }}</div>
            <div v-else style="margin-top:12px; flex:1; min-height:0; overflow-y:auto;">
              <div v-if="formulaPackages.length > 0">
                <div style="font-size:11px; font-weight:600; text-transform:uppercase; letter-spacing:0.5px; color:var(--color-text-tertiary); margin-bottom:6px;">{{ t('taps.formulaeCount', { count: formulaPackages.length }) }}</div>
                <div class="card-group" style="margin-bottom:12px;">
                  <div v-for="pkg in formulaPackages" :key="pkg.id" class="card-row">
                    <div style="display:flex; align-items:center; gap:8px; min-width:0; flex:1;">
                      <input type="checkbox" style="flex-shrink:0;" :checked="isSelectedPkg(pkg.id)" @change="togglePkg(pkg.id)" />
                      <PackageIcon :name="pkg.name" :size="28" />
                      <span style="font-size:13px; color:var(--color-text); overflow:hidden; text-overflow:ellipsis; white-space:nowrap;">{{ pkg.name }}</span>
                      <span v-if="pkg.installed" class="pkg-badge installed">{{ t('common.installed') }}</span>
                      <span v-if="pkg.outdated" class="pkg-badge outdated">{{ t('common.outdated') }}</span>
                    </div>
                    <div style="display:flex; align-items:center; gap:2px;">
                      <button class="icon-btn accent" :disabled="pkg.installed" @click="runPackageAction(pkg, 'install')"><Download :size="14" /></button>
                      <button class="icon-btn" style="color:var(--color-warning);" :disabled="!pkg.installed || !pkg.outdated" @click="runPackageAction(pkg, 'upgrade')"><ArrowUpCircle :size="14" /></button>
                      <button class="icon-btn danger" :disabled="!pkg.installed" @click="runPackageAction(pkg, 'uninstall')"><Trash2 :size="14" /></button>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="caskPackages.length > 0">
                <div style="font-size:11px; font-weight:600; text-transform:uppercase; letter-spacing:0.5px; color:var(--color-text-tertiary); margin-bottom:6px;">{{ t('taps.casksCount', { count: caskPackages.length }) }}</div>
                <div class="card-group">
                  <div v-for="pkg in caskPackages" :key="pkg.id" class="card-row">
                    <div style="display:flex; align-items:center; gap:8px; min-width:0; flex:1;">
                      <input type="checkbox" style="flex-shrink:0;" :checked="isSelectedPkg(pkg.id)" @change="togglePkg(pkg.id)" />
                      <PackageIcon :name="pkg.name" :size="28" />
                      <span style="font-size:13px; color:var(--color-text); overflow:hidden; text-overflow:ellipsis; white-space:nowrap;">{{ pkg.name }}</span>
                      <span v-if="pkg.installed" class="pkg-badge installed">{{ t('common.installed') }}</span>
                      <span v-if="pkg.outdated" class="pkg-badge outdated">{{ t('common.outdated') }}</span>
                    </div>
                    <div style="display:flex; align-items:center; gap:2px;">
                      <button class="icon-btn accent" :disabled="pkg.installed" @click="runPackageAction(pkg, 'install')"><Download :size="14" /></button>
                      <button class="icon-btn" style="color:var(--color-warning);" :disabled="!pkg.installed || !pkg.outdated" @click="runPackageAction(pkg, 'upgrade')"><ArrowUpCircle :size="14" /></button>
                      <button class="icon-btn danger" :disabled="!pkg.installed" @click="runPackageAction(pkg, 'uninstall')"><Trash2 :size="14" /></button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </teleport>

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
.icon-btn.danger {
  color: var(--color-danger);
}
.icon-btn.danger:hover {
  background: var(--color-danger-light);
}
.icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 40;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.4);
  padding: 16px;
}
.modal-content {
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  width: 100%;
  max-width: 800px;
  border-radius: var(--radius-lg);
  background: var(--color-card);
  padding: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.pkg-badge {
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}
.pkg-badge.installed {
  background: var(--color-success-light);
  color: var(--color-success);
}
.pkg-badge.outdated {
  background: var(--color-warning-light);
  color: var(--color-warning);
}
</style>
