<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, ArrowUpCircle, CheckCircle } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useUpdateStore } from '@/stores/update'
import { useLogStore } from '@/stores/log'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Toast from '@/components/common/Toast.vue'
import BulkActionBar from '@/components/common/BulkActionBar.vue'
import BulkResultSummary from '@/components/common/BulkResultSummary.vue'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'
import PackageIcon from '@/components/icons/PackageIcon.vue'

const { t } = useI18n()
const updateStore = useUpdateStore()
const logStore = useLogStore()
const selection = useSelection()
const toastRef = ref<any>(null)
const upgradeAllOpen = ref(false)
const bulkSummary = ref<BulkSummary | null>(null)
const updatingBrew = ref(false)

function sortByName<T extends { name: string }>(items: T[]) {
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
}

const formulae = computed(() => sortByName(updateStore.formulae))
const casks = computed(() => sortByName(updateStore.casks))
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

async function upgrade(name: string) {
  logStore.startListening()
  try {
    await updateStore.upgrade(name)
    toastRef.value?.show('success', t('messages.upgradedName', { name }))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.upgradeFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function upgradeSelected() {
  const names = selectedNames.value
  if (names.length === 0) return

  logStore.startListening()
  try {
    const result = await runBulk(names, (name) => BrewService.Upgrade(name))
    bulkSummary.value = {
      action: t('messages.bulkUpgrade'),
      total: result.total,
      success: result.success,
      failures: result.failures,
      timestamp: Date.now(),
    }
    await updateStore.fetchOutdated()
    selection.clear()
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', t('messages.bulkUpgradedSummary', { success: result.success, total: result.total }))
  } finally {
    logStore.stopListening()
  }
}

async function upgradeAll() {
  upgradeAllOpen.value = false
  logStore.startListening()
  try {
    await updateStore.upgradeAll()
    toastRef.value?.show('success', t('messages.allPackagesUpgraded'))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.upgradeFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function runBrewUpdate() {
  if (updatingBrew.value) return
  updatingBrew.value = true
  logStore.startListening()
  try {
    await BrewService.Update()
    await updateStore.fetchOutdated()
    toastRef.value?.show('success', t('messages.brewUpdateCompleted'))
  } catch (error: any) {
    const code = error?.cause?.code || error?.code
    if (code === 'BREW_BUSY') {
      toastRef.value?.show('error', t('messages.brewUpdateBusy'))
    } else {
      toastRef.value?.show('error', error?.message || t('messages.brewUpdateFailed'))
    }
  } finally {
    logStore.stopListening()
    updatingBrew.value = false
  }
}

onMounted(() => {
  updateStore.fetchOutdated()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header" style="display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px;">
      <h1 class="content-title">{{ t('update.title') }}</h1>
      <div class="action-bar">
        <button
          style="background:transparent; color:var(--color-text); border:1px solid var(--color-border); border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; cursor:pointer; display:inline-flex; align-items:center; gap:4px;"
          :disabled="updatingBrew || updateStore.updating"
          @click="runBrewUpdate"
        >
          <RefreshCw :size="14" />
          {{ t('update.brewUpdate') }}
        </button>
        <button
          style="background:transparent; border:none; border-radius:var(--radius-sm); padding:6px; cursor:pointer; display:inline-flex; align-items:center; color:var(--color-text-tertiary);"
          @click="updateStore.forceRefresh()"
        >
          <RefreshCw :size="14" />
        </button>
        <button
          style="background:var(--color-accent); color:white; border:none; border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; font-weight:500; cursor:pointer;"
          :disabled="updateStore.updating || (formulae.length === 0 && casks.length === 0)"
          @click="upgradeAllOpen = true"
        >
          {{ t('update.upgradeAll') }}
        </button>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <LoadingSkeleton v-if="updateStore.loading" :rows="4" />

      <div v-else-if="updateStore.error" style="padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">
        {{ updateStore.error }}
      </div>

      <div v-else-if="formulae.length === 0 && casks.length === 0" style="text-align:center; padding:40px 0;">
        <CheckCircle :size="40" style="color:var(--color-success); margin:0 auto 12px;" />
        <p style="font-size:13px; color:var(--color-success);">{{ t('update.upToDate') }}</p>
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
            @click="upgradeSelected"
          >
            {{ t('update.upgrade') }}
          </button>
        </BulkActionBar>

        <div v-if="formulae.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.formulae') }}</h3>
          <div class="card-group">
            <div v-for="item in formulae" :key="`formula-${item.name}`" class="card-row">
              <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
                <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(`f:${item.name}`)" @change="selection.toggle(`f:${item.name}`)" />
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.installed_versions?.[0] || '' }} → {{ item.current_version }}</div>
                </div>
              </div>
              <button class="icon-btn accent" @click="upgrade(item.name)">
                <ArrowUpCircle :size="16" />
              </button>
            </div>
          </div>
        </div>

        <div v-if="casks.length > 0" class="content-section">
          <h3 class="section-title">{{ t('installed.casks') }}</h3>
          <div class="card-group">
            <div v-for="item in casks" :key="`cask-${item.name}`" class="card-row">
              <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
                <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(`c:${item.name}`)" @change="selection.toggle(`c:${item.name}`)" />
                <PackageIcon :name="item.name" :size="28" />
                <div style="min-width:0;">
                  <div class="card-row-title">{{ item.name }}</div>
                  <div class="card-row-subtitle">{{ item.installed_version || '' }} → {{ item.current_version }}</div>
                </div>
              </div>
              <button class="icon-btn accent" @click="upgrade(item.name)">
                <ArrowUpCircle :size="16" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :open="upgradeAllOpen"
      :title="t('update.upgradeAll')"
      :message="t('update.confirmUpgradeAll')"
      :confirm-label="t('update.upgradeAll')"
      @confirm="upgradeAll"
      @cancel="upgradeAllOpen = false"
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
}
.icon-btn.accent {
  color: var(--color-accent);
}
.icon-btn.accent:hover {
  background: var(--color-accent-light);
}
</style>
