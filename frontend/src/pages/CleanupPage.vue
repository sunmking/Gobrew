<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Trash2, Loader2 } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useLogStore } from '@/stores/log'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import LoadingInline from '@/components/common/LoadingInline.vue'
import LoadingSkeleton from '@/components/common/LoadingSkeleton.vue'
import Toast from '@/components/common/Toast.vue'

const props = withDefaults(
  defineProps<{
    embedded?: boolean
  }>(),
  {
    embedded: false,
  },
)

const { t } = useI18n()
const logStore = useLogStore()
const toastRef = ref<any>(null)

const activeTab = ref<'cleanup' | 'autoremove'>('cleanup')
const preview = ref('')
const autoremovePreview = ref('')
const loading = ref(false)
const autoremoveLoading = ref(false)
const error = ref('')
const autoremoveError = ref('')
const showCleanupDialog = ref(false)
const showAutoremoveDialog = ref(false)
const query = ref('')
const autoremoveQuery = ref('')
const refreshing = ref(false)
const executingCleanup = ref(false)
const executingAutoremove = ref(false)

const previewLines = computed(() => {
  const lines = (preview.value || '')
    .split('\n')
    .map((line) => line.trimEnd())
    .filter((line) => line.trim() !== '')
    .sort((a, b) => a.localeCompare(b))
  if (!query.value.trim()) return lines
  const q = query.value.trim().toLowerCase()
  return lines.filter((line) => line.toLowerCase().includes(q))
})

const autoremoveLines = computed(() => {
  const lines = (autoremovePreview.value || '')
    .split('\n')
    .map((line) => line.trimEnd())
    .filter((line) => line.trim() !== '')
    .sort((a, b) => a.localeCompare(b))
  if (!autoremoveQuery.value.trim()) return lines
  const q = autoremoveQuery.value.trim().toLowerCase()
  return lines.filter((line) => line.toLowerCase().includes(q))
})

async function fetchPreview() {
  loading.value = true
  error.value = ''
  try {
    const result = await BrewService.CleanupPreview()
    preview.value = result?.output || ''
  } catch (e: any) {
    preview.value = ''
    error.value = e?.message || t('messages.failedLoadCleanupPreview')
  } finally {
    loading.value = false
  }
}

async function fetchAutoremovePreview() {
  autoremoveLoading.value = true
  autoremoveError.value = ''
  try {
    const result = await BrewService.AutoRemovePreview()
    autoremovePreview.value = result?.output || ''
  } catch (e: any) {
    autoremovePreview.value = ''
    autoremoveError.value = e?.message || t('messages.failedLoadAutoremovePreview')
  } finally {
    autoremoveLoading.value = false
  }
}

async function refreshAll() {
  refreshing.value = true
  try {
    await Promise.all([fetchPreview(), fetchAutoremovePreview()])
  } finally {
    refreshing.value = false
  }
}

async function executeCleanup() {
  showCleanupDialog.value = false
  executingCleanup.value = true
  logStore.startListening()
  try {
    await BrewService.Cleanup()
    toastRef.value?.show('success', t('messages.cleanupComplete'))
    await fetchPreview()
  } catch (e: any) {
    toastRef.value?.show('error', e?.message || t('messages.cleanupFailed'))
  } finally {
    logStore.stopListening()
    executingCleanup.value = false
  }
}

async function executeAutoremove() {
  showAutoremoveDialog.value = false
  executingAutoremove.value = true
  logStore.startListening()
  try {
    await BrewService.AutoRemove()
    toastRef.value?.show('success', t('cleanup.autoremoveDone'))
    await fetchAutoremovePreview()
  } catch (e: any) {
    toastRef.value?.show('error', e?.message || t('messages.autoremoveFailed'))
  } finally {
    logStore.stopListening()
    executingAutoremove.value = false
  }
}

onMounted(async () => {
  await refreshAll()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header" :style="embedded ? 'margin-bottom:10px;' : ''">
      <div style="display:flex; align-items:center; justify-content:space-between;">
        <h1 v-if="!embedded" class="content-title">{{ t('cleanup.title') }}</h1>
        <div v-else />
        <button
          class="ui-btn ui-btn-secondary ui-btn-sm"
          :disabled="refreshing || loading || autoremoveLoading"
          @click="refreshAll"
        >
          <RefreshCw :size="14" :class="refreshing ? 'is-spinning' : ''" />
          {{ t('common.refresh') }}
        </button>
      </div>
      <div class="ui-segmented" style="margin-top:12px;">
        <button
          class="ui-segmented-btn"
          :class="{ 'is-active': activeTab === 'cleanup' }"
          @click="activeTab = 'cleanup'"
        >
          Cleanup
        </button>
        <button
          class="ui-segmented-btn"
          :class="{ 'is-active': activeTab === 'autoremove' }"
          @click="activeTab = 'autoremove'"
        >
          Autoremove
        </button>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <div v-if="activeTab === 'cleanup'">
        <div style="display:flex; align-items:center; justify-content:space-between; gap:12px;">
          <h3 class="section-title">{{ t('cleanup.preview') }}</h3>
          <button class="ui-btn ui-btn-danger ui-btn-sm" :disabled="executingCleanup || loading" @click="showCleanupDialog = true">
            <Trash2 v-if="!executingCleanup" :size="14" />
            <Loader2 v-else :size="14" class="ui-spinner" />
            {{ executingCleanup ? t('common.loading') : t('cleanup.execute') }}
          </button>
        </div>
        <input
          v-model="query"
          type="text"
          :placeholder="t('placeholders.filterPreview')"
          class="ui-input"
          style="width:100%; margin-top:8px;"
        />
        <div v-if="loading" style="margin-top:12px;">
          <LoadingSkeleton :rows="4" :show-header="false" />
        </div>
        <div v-else-if="error" style="margin-top:12px; padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">{{ error }}</div>
        <div v-else-if="previewLines.length > 0" class="card-group" style="margin-top:12px;">
          <div v-for="(line, i) in previewLines" :key="i" class="card-row" style="font-family:var(--font-mono); font-size:12px;">
            <span style="color:var(--color-text-tertiary); width:28px; text-align:right; flex-shrink:0; user-select:none;">{{ i + 1 }}</span>
            <span style="color:var(--color-text); word-break:break-all;">{{ line }}</span>
          </div>
        </div>
        <div v-else style="margin-top:12px; font-size:13px; color:var(--color-text-tertiary);">{{ t('cleanup.noCleanup') }}</div>
      </div>

      <div v-if="activeTab === 'autoremove'">
        <div style="display:flex; align-items:center; justify-content:space-between; gap:12px;">
          <h3 class="section-title">{{ t('cleanup.autoremovePreview') }}</h3>
          <button class="ui-btn ui-btn-primary ui-btn-sm" :disabled="executingAutoremove || autoremoveLoading" @click="showAutoremoveDialog = true">
            <Trash2 v-if="!executingAutoremove" :size="14" />
            <Loader2 v-else :size="14" class="ui-spinner" />
            {{ executingAutoremove ? t('common.loading') : t('cleanup.autoremove') }}
          </button>
        </div>
        <input
          v-model="autoremoveQuery"
          type="text"
          :placeholder="t('cleanup.filterAutoremove')"
          class="ui-input"
          style="width:100%; margin-top:8px;"
        />
        <div v-if="autoremoveLoading" style="margin-top:12px;">
          <LoadingInline />
        </div>
        <div v-else-if="autoremoveError" style="margin-top:12px; padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">{{ autoremoveError }}</div>
        <div v-else-if="autoremoveLines.length > 0" class="card-group" style="margin-top:12px;">
          <div v-for="(line, i) in autoremoveLines" :key="i" class="card-row" style="font-family:var(--font-mono); font-size:12px;">
            <span style="color:var(--color-text-tertiary); width:28px; text-align:right; flex-shrink:0; user-select:none;">{{ i + 1 }}</span>
            <span style="color:var(--color-text); word-break:break-all;">{{ line }}</span>
          </div>
        </div>
        <div v-else style="margin-top:12px; font-size:13px; color:var(--color-text-tertiary);">{{ t('cleanup.noAutoremove') }}</div>
      </div>
    </div>

    <ConfirmDialog
      :open="showCleanupDialog"
      :title="t('cleanup.execute')"
      :message="t('cleanup.confirmExecute')"
      :confirm-label="t('cleanup.execute')"
      severity="danger"
      @confirm="executeCleanup"
      @cancel="showCleanupDialog = false"
    />
    <ConfirmDialog
      :open="showAutoremoveDialog"
      :title="t('cleanup.autoremove')"
      :message="t('cleanup.confirmAutoremove')"
      :confirm-label="t('cleanup.autoremove')"
      severity="danger"
      @confirm="executeAutoremove"
      @cancel="showAutoremoveDialog = false"
    />
    <Toast ref="toastRef" />
  </div>
</template>
