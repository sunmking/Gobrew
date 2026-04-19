<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Save, Download, CheckCircle, RefreshCw, Wrench, Loader2 } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useBundleStore } from '@/stores/bundle'
import { useLogStore } from '@/stores/log'
import { useSettingsStore } from '@/stores/settings'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import LoadingInline from '@/components/common/LoadingInline.vue'
import Toast from '@/components/common/Toast.vue'
import BulkActionBar from '@/components/common/BulkActionBar.vue'
import BulkResultSummary from '@/components/common/BulkResultSummary.vue'

const props = withDefaults(
  defineProps<{
    embedded?: boolean
  }>(),
  {
    embedded: false,
  },
)

const { t } = useI18n()
const bundleStore = useBundleStore()
const logStore = useLogStore()
const settingsStore = useSettingsStore()
const selection = useSelection()
const toastRef = ref<any>(null)

const showRestoreDialog = ref(false)
const showCleanupDialog = ref(false)
const cleanupForce = ref(false)
const query = ref('')
const bulkSummary = ref<BulkSummary | null>(null)
const loadingBrewfile = ref(false)
const dumping = ref(false)
const restoring = ref(false)
const checking = ref(false)
const cleanupPreviewLoading = ref(false)
const cleaningUp = ref(false)
const installingMissing = ref(false)

const brewfileLines = computed(() => {
  const lines = (bundleStore.brewfileContent || '')
    .split('\n')
    .map((line) => line.trimEnd())
    .filter((line) => line.trim() !== '')
  if (!query.value.trim()) return lines
  const q = query.value.trim().toLowerCase()
  return lines.filter((line) => line.toLowerCase().includes(q))
})

function lineColor(line: string): string {
  const trimmed = line.trim()
  if (trimmed.startsWith('#')) return 'var(--color-text-tertiary)'
  if (trimmed.startsWith('tap ')) return 'var(--color-accent)'
  if (trimmed.startsWith('brew ')) return 'var(--color-success)'
  if (trimmed.startsWith('cask ')) return 'var(--color-warning)'
  return 'var(--color-success)'
}

const missingItems = computed(() => [...(bundleStore.checkResult?.missing || [])].sort((a, b) => a.localeCompare(b)))
const allSelected = computed(() => selection.isAllSelected(missingItems.value))

watch(
  missingItems,
  (ids) => {
    selection.sync(ids)
  },
  { immediate: true },
)

async function loadBrewfile() {
  loadingBrewfile.value = true
  try {
    await bundleStore.readBrewfile(settingsStore.brewFilePath || undefined)
  } finally {
    loadingBrewfile.value = false
  }
}

async function dump() {
  dumping.value = true
  logStore.startListening()
  try {
    await bundleStore.dump(settingsStore.brewFilePath || undefined, true)
    toastRef.value?.show('success', `${t('bundle.brewfile')} dumped`)
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.dumpFailed'))
  } finally {
    logStore.stopListening()
    dumping.value = false
  }
}

async function restore() {
  showRestoreDialog.value = false
  restoring.value = true
  logStore.startListening()
  try {
    await bundleStore.restore(settingsStore.brewFilePath || undefined)
    toastRef.value?.show('success', `${t('bundle.brewfile')} restored`)
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.restoreFailed'))
  } finally {
    logStore.stopListening()
    restoring.value = false
  }
}

async function check() {
  checking.value = true
  try {
    await bundleStore.check(settingsStore.brewFilePath || undefined)
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.checkFailed'))
  } finally {
    checking.value = false
  }
}

async function cleanupPreview() {
  cleanupPreviewLoading.value = true
  try {
    await bundleStore.cleanupPreview(settingsStore.brewFilePath || undefined)
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.cleanupPreviewFailed'))
  } finally {
    cleanupPreviewLoading.value = false
  }
}

async function cleanup() {
  showCleanupDialog.value = false
  cleaningUp.value = true
  logStore.startListening()
  try {
    await bundleStore.cleanup(settingsStore.brewFilePath || undefined, cleanupForce.value)
    toastRef.value?.show('success', t('bundle.cleanupDone'))
    await Promise.all([check(), cleanupPreview()])
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.cleanupFailed'))
  } finally {
    logStore.stopListening()
    cleaningUp.value = false
  }
}

async function installSelectedMissing() {
  const names = [...selection.selected.value]
  if (names.length === 0) return

  installingMissing.value = true
  logStore.startListening()
  try {
    const result = await runBulk(names, (name) => BrewService.Install(name))
    bulkSummary.value = {
      action: 'Install Missing from Brewfile',
      total: result.total,
      success: result.success,
      failures: result.failures,
      timestamp: Date.now(),
    }
    await check()
    selection.clear()
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', t('messages.bulkInstalledSummary', { success: result.success, total: result.total }))
  } finally {
    logStore.stopListening()
    installingMissing.value = false
  }
}

onMounted(async () => {
  await loadBrewfile()
  await cleanupPreview()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header" :style="embedded ? 'margin-bottom:10px;' : ''">
      <div style="display:flex; align-items:center; justify-content:space-between;">
        <h1 v-if="!embedded" class="content-title">{{ t('bundle.title') }}</h1>
        <div v-else />
        <button
          class="ui-btn ui-btn-secondary ui-btn-sm"
          :disabled="loadingBrewfile"
          @click="loadBrewfile"
        >
          <RefreshCw :size="14" :class="loadingBrewfile ? 'is-spinning' : ''" />
          {{ t('common.refresh') }}
        </button>
      </div>
      <div style="display:flex; flex-direction:column; gap:8px; margin-top:12px;">
        <div style="display:flex; gap:8px; flex-wrap:wrap;">
          <input
            v-model="settingsStore.brewFilePath"
            type="text"
            :placeholder="t('bundle.brewfilePath')"
            class="ui-input"
            style="flex:1; min-width:200px;"
            @change="settingsStore.setBrewFilePath(settingsStore.brewFilePath); loadBrewfile()"
          />
          <input
            v-model="query"
            type="text"
            :placeholder="t('placeholders.filterBrewfile')"
            class="ui-input"
            style="width:180px;"
          />
        </div>
        <div class="action-bar">
          <button
            class="ui-btn ui-btn-primary ui-btn-sm"
            :class="{ 'is-loading': dumping }"
            :disabled="dumping"
            @click="dump"
          >
            <Save v-if="!dumping" :size="14" />
            <Loader2 v-else :size="14" class="ui-spinner" />
            {{ dumping ? t('common.loading') : t('bundle.dump') }}
          </button>
          <button
            class="ui-btn ui-btn-primary ui-btn-sm"
            :disabled="restoring"
            @click="showRestoreDialog = true"
          >
            <Download :size="14" />
            {{ t('bundle.restore') }}
          </button>
          <button
            class="ui-btn ui-btn-secondary ui-btn-sm"
            :disabled="checking"
            @click="check"
          >
            <CheckCircle :size="14" :class="checking ? 'is-spinning' : ''" />
            {{ t('bundle.check') }}
          </button>
          <button
            class="ui-btn ui-btn-danger ui-btn-sm"
            :disabled="cleaningUp"
            @click="showCleanupDialog = true"
          >
            <Wrench :size="14" />
            {{ t('bundle.cleanup') }}
          </button>
        </div>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <div v-if="bundleStore.error" style="margin-bottom:12px; padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">
        {{ bundleStore.error }}
      </div>

      <div v-if="brewfileLines.length > 0" class="card-group">
        <div v-for="(line, i) in brewfileLines" :key="i" class="card-row" style="font-family:var(--font-mono); font-size:12px;">
          <span style="color:var(--color-text-tertiary); width:28px; text-align:right; flex-shrink:0; user-select:none;">{{ i + 1 }}</span>
          <span :style="`color:${lineColor(line)}; word-break:break-all;`">{{ line }}</span>
        </div>
      </div>
      <div v-else style="text-align:center; padding:32px 0; font-size:13px; color:var(--color-text-tertiary);">
        {{ t('bundle.noBrewfile') }}
      </div>

      <div v-if="bundleStore.checkResult" style="margin-top:16px;">
        <div v-if="bundleStore.checkResult.satisfied" style="padding:12px; border-radius:var(--radius-sm); background:var(--color-success-light); color:var(--color-success); font-size:13px;">
          {{ t('bundle.allUpToDate') }}
        </div>
        <div v-else>
          <div style="padding:12px; border-radius:var(--radius-sm); background:var(--color-accent-light); margin-bottom:12px;">
            <h3 style="font-size:13px; font-weight:500; color:var(--color-accent); margin:0 0 8px;">{{ t('bundle.missing') }}</h3>

            <BulkActionBar
              :selected-count="selection.selectedCount.value"
              :all-selected="allSelected"
              @select-all="selection.selectAll(missingItems)"
              @clear="selection.clear()"
            >
              <button
                class="btn-primary"
                style="display:inline-flex; align-items:center; gap:4px;"
                :class="{ 'is-loading': installingMissing }"
                :disabled="selection.selectedCount.value === 0 || installingMissing"
                @click="installSelectedMissing"
              >
                <Wrench v-if="!installingMissing" :size="12" />
                <Loader2 v-else :size="12" class="ui-spinner" />
                {{ installingMissing ? t('common.loading') : t('bundle.installSelected') }}
              </button>
            </BulkActionBar>

            <div class="card-group" style="margin-top:8px;">
              <div v-for="item in missingItems" :key="item" class="card-row">
                <div style="display:flex; align-items:center; gap:8px; min-width:0; flex:1;">
                  <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(item)" @change="selection.toggle(item)" />
                  <span style="font-size:13px; color:var(--color-accent);">{{ item }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="bundleStore.cleanupPreviewResult?.output || bundleStore.cleanupPreviewResult === null" style="margin-top:16px;">
        <div style="display:flex; align-items:center; justify-content:space-between; gap:8px; margin-bottom:8px;">
          <h3 class="section-title" style="margin:0;">{{ t('bundle.cleanupPreview') }}</h3>
          <button
            class="ui-btn ui-btn-secondary ui-btn-sm"
            :disabled="cleanupPreviewLoading"
            @click="cleanupPreview"
          >
            <RefreshCw :size="14" :class="cleanupPreviewLoading ? 'is-spinning' : ''" />
            {{ t('common.refresh') }}
          </button>
        </div>
        <div
          v-if="bundleStore.cleanupPreviewResult?.output"
          style="padding:12px; border-radius:var(--radius-sm); background:var(--color-card); border:1px solid var(--color-border); font-family:var(--font-mono); font-size:12px; color:var(--color-text); white-space:pre-wrap;"
        >
          {{ bundleStore.cleanupPreviewResult.output }}
        </div>
        <div v-else style="font-size:13px; color:var(--color-text-tertiary);">
          <LoadingInline v-if="cleanupPreviewLoading" />
          <span v-else>{{ t('bundle.noCleanupPreview') }}</span>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :open="showRestoreDialog"
      :title="t('bundle.restore')"
      :message="t('bundle.confirmRestore')"
      :confirm-label="t('bundle.restore')"
      @confirm="restore"
      @cancel="showRestoreDialog = false"
    />
    <ConfirmDialog
      :open="showCleanupDialog"
      :title="t('bundle.cleanup')"
      :message="t('bundle.confirmCleanup')"
      :confirm-label="t('bundle.cleanup')"
      severity="danger"
      :toggle-label="t('bundle.forceCleanup')"
      :toggle-value="cleanupForce"
      @update:toggle-value="(value) => (cleanupForce = value)"
      @confirm="cleanup"
      @cancel="showCleanupDialog = false"
    />
    <Toast ref="toastRef" />
  </div>
</template>
