<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Play, Square, RotateCw } from 'lucide-vue-next'
import * as ServiceManager from '../../bindings/changeme/services/servicemanager.js'
import { useServicesStore } from '@/stores/services'
import { useLogStore } from '@/stores/log'
import { useSelection } from '@/composables/useSelection'
import { runBulk } from '@/composables/useBulkRunner'
import type { BulkSummary } from '@/types/bulk'
import Toast from '@/components/common/Toast.vue'
import BulkActionBar from '@/components/common/BulkActionBar.vue'
import BulkResultSummary from '@/components/common/BulkResultSummary.vue'
import PackageIcon from '@/components/icons/PackageIcon.vue'

const props = withDefaults(
  defineProps<{
    embedded?: boolean
  }>(),
  {
    embedded: false,
  },
)

const { t } = useI18n()
const servicesStore = useServicesStore()
const logStore = useLogStore()
const selection = useSelection()
const toastRef = ref<any>(null)
const query = ref('')
const bulkSummary = ref<BulkSummary | null>(null)

function sortByName<T extends { name: string }>(items: T[]) {
  return [...items].sort((a, b) => a.name.localeCompare(b.name))
}

const services = computed(() => {
  const q = query.value.trim().toLowerCase()
  const items = q
    ? servicesStore.services.filter((item) => [item.name, item.status, item.user, item.file].join(' ').toLowerCase().includes(q))
    : servicesStore.services
  return sortByName(items)
})

const allVisibleKeys = computed(() => services.value.map((item) => item.name))
const allSelected = computed(() => selection.isAllSelected(allVisibleKeys.value))

watch(
  allVisibleKeys,
  (ids) => {
    selection.sync(ids)
  },
  { immediate: true },
)

async function runAction(action: 'start' | 'stop' | 'restart', name: string) {
  logStore.startListening()
  try {
    if (action === 'start') await servicesStore.start(name)
    if (action === 'stop') await servicesStore.stop(name)
    if (action === 'restart') await servicesStore.restart(name)
    toastRef.value?.show('success', t('messages.serviceActionDone', { name, action }))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.operationFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function runSelected(action: 'start' | 'stop' | 'restart') {
  const names = [...selection.selected.value]
  if (names.length === 0) return

  logStore.startListening()
  try {
    const result = await runBulk(names, async (name) => {
      if (action === 'start') await ServiceManager.Start(name)
      if (action === 'stop') await ServiceManager.Stop(name)
      if (action === 'restart') await ServiceManager.Restart(name)
    })
    bulkSummary.value = {
      action: t('messages.bulkActionLabel', { action }),
      total: result.total,
      success: result.success,
      failures: result.failures,
      timestamp: Date.now(),
    }
    await servicesStore.fetchServices()
    selection.clear()
    toastRef.value?.show(result.failures.length === 0 ? 'success' : 'error', t('messages.bulkProcessedSummary', { success: result.success, total: result.total }))
  } finally {
    logStore.stopListening()
  }
}

async function runAll(action: 'start' | 'stop' | 'restart') {
  logStore.startListening()
  try {
    if (action === 'start') await servicesStore.startAll()
    if (action === 'stop') await servicesStore.stopAll()
    if (action === 'restart') await servicesStore.restartAll()
    toastRef.value?.show('success', t(`services.${action}All`))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.operationFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function runCleanup() {
  logStore.startListening()
  try {
    await servicesStore.cleanup()
    toastRef.value?.show('success', t('services.cleanupDone'))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.operationFailed'))
  } finally {
    logStore.stopListening()
  }
}

onMounted(() => {
  servicesStore.fetchServices()
})
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <div class="content-header">
      <div style="display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:8px;">
        <h1 class="content-title">{{ t('services.title') }}</h1>
        <div class="action-bar">
          <button
            style="background:var(--color-success); color:white; border:none; border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; font-weight:500; cursor:pointer;"
            @click="runAll('start')"
          >
            {{ t('services.startAll') }}
          </button>
          <button
            style="background:var(--color-danger); color:white; border:none; border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; font-weight:500; cursor:pointer;"
            @click="runAll('stop')"
          >
            {{ t('services.stopAll') }}
          </button>
          <button
            style="background:var(--color-accent); color:white; border:none; border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; font-weight:500; cursor:pointer;"
            @click="runAll('restart')"
          >
            {{ t('services.restartAll') }}
          </button>
          <button
            style="background:transparent; color:var(--color-text-secondary); border:1px solid var(--color-border); border-radius:var(--radius-sm); padding:6px 14px; font-size:13px; cursor:pointer;"
            @click="runCleanup"
          >
            {{ t('services.cleanup') }}
          </button>
          <button
            style="background:transparent; border:none; border-radius:var(--radius-sm); padding:6px; cursor:pointer; display:inline-flex; align-items:center; color:var(--color-text-tertiary);"
            @click="servicesStore.fetchServices()"
          >
            <RefreshCw :size="14" />
          </button>
        </div>
      </div>
      <input
        v-model="query"
        :placeholder="t('placeholders.searchServices')"
        style="width:100%; margin-top:12px; padding:8px 12px; font-size:13px; border-radius:var(--radius-sm); border:1px solid var(--color-border); background:var(--color-card); color:var(--color-text); outline:none;"
      />
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto">
      <BulkResultSummary :summary="bulkSummary" />

      <div v-if="servicesStore.loading" style="font-size:13px; color:var(--color-text-tertiary);">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="servicesStore.error" style="padding:12px; border-radius:var(--radius-sm); background:var(--color-danger-light); color:var(--color-danger); font-size:13px;">
        {{ servicesStore.error }}
      </div>

      <div v-else-if="services.length === 0" style="text-align:center; padding:40px 0; font-size:13px; color:var(--color-text-tertiary);">
        {{ t('services.noServices') }}
      </div>

      <div v-else>
        <BulkActionBar
          :selected-count="selection.selectedCount.value"
          :all-selected="allSelected"
          @select-all="selection.selectAll(allVisibleKeys)"
          @clear="selection.clear()"
        >
          <button class="btn-primary" :disabled="selection.selectedCount.value === 0" @click="runSelected('start')">{{ t('services.start') }}</button>
          <button class="btn-danger" :disabled="selection.selectedCount.value === 0" @click="runSelected('stop')">{{ t('services.stop') }}</button>
          <button class="btn-neutral" :disabled="selection.selectedCount.value === 0" @click="runSelected('restart')">{{ t('services.restart') }}</button>
        </BulkActionBar>

        <div class="card-group">
          <div v-for="service in services" :key="service.name" class="card-row">
            <div style="display:flex; align-items:center; gap:10px; min-width:0; flex:1;">
              <input type="checkbox" style="flex-shrink:0;" :checked="selection.isSelected(service.name)" @change="selection.toggle(service.name)" />
              <PackageIcon :name="service.name" :size="28" />
              <div style="min-width:0;">
                <div style="display:flex; align-items:center; gap:8px;">
                  <span class="card-row-title">{{ service.name }}</span>
                  <span
                    class="status-badge"
                    :class="{
                      'status-started': service.status === 'started',
                      'status-stopped': service.status === 'stopped' || service.status === 'none',
                      'status-error': service.status === 'error',
                    }"
                  >
                    <span class="status-dot" :class="{
                      'dot-started': service.status === 'started',
                      'dot-stopped': service.status === 'stopped' || service.status === 'none',
                      'dot-error': service.status === 'error',
                    }" />
                    {{ service.status }}
                  </span>
                </div>
                <div class="card-row-subtitle">{{ service.user }}</div>
                <div v-if="service.file" style="font-size:12px; color:var(--color-text-tertiary); margin-top:2px;">{{ service.file }}</div>
              </div>
            </div>
            <div style="display:flex; align-items:center; gap:2px;">
              <button class="icon-btn success" @click="runAction('start', service.name)"><Play :size="15" /></button>
              <button class="icon-btn danger" @click="runAction('stop', service.name)"><Square :size="15" /></button>
              <button class="icon-btn accent" @click="runAction('restart', service.name)"><RotateCw :size="15" /></button>
            </div>
          </div>
        </div>
      </div>
    </div>

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

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}
.status-badge.status-started {
  background: var(--color-success-light);
  color: var(--color-success);
}
.status-badge.status-stopped {
  background: var(--color-card);
  color: var(--color-text-tertiary);
}
.status-badge.status-error {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.status-dot.dot-started {
  background: var(--color-success);
}
.status-dot.dot-stopped {
  background: var(--color-text-tertiary);
}
.status-dot.dot-error {
  background: var(--color-danger);
}
</style>
