<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Trash2, Upload } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import BrewButton from '@/components/common/BrewButton.vue'
import SegmentedControl from '@/components/common/SegmentedControl.vue'
import StatCard from '@/components/common/StatCard.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import TerminalPanel from '@/components/common/TerminalPanel.vue'
import { useUpdateStore } from '@/stores/update'
import { useLogStore } from '@/stores/log'

const updateStore = useUpdateStore()
const logStore = useLogStore()
const { t } = useI18n()
const tab = ref('upgrades')
const cleanupOutput = ref('')
const autoremoveOutput = ref('')
const error = ref('')

const totalOutdated = computed(() => updateStore.formulae.length + updateStore.casks.length)
const totalPinned = computed(() => updateStore.formulae.filter(item => item.pinned).length)
const tabs = computed(() => [
  { label: t('updateCleanup.upgradesTab'), value: 'upgrades' },
  { label: t('updateCleanup.cleanupTab'), value: 'cleanup' },
  { label: t('updateCleanup.cacheTab'), value: 'cache' },
])

async function refresh() {
  error.value = ''
  try {
    await updateStore.forceRefresh()
    const [cleanup, autoremove] = await Promise.all([BrewService.CleanupPreview(), BrewService.AutoRemovePreview()])
    cleanupOutput.value = cleanup?.output || ''
    autoremoveOutput.value = autoremove?.output || ''
  } catch (err: any) {
    error.value = err?.message || t('messages.operationFailed')
  }
}

async function run(name: string, fn: () => Promise<void>) {
  logStore.startListening(name)
  try {
    await fn()
    await refresh()
  } finally {
    logStore.stopListening()
  }
}

onMounted(refresh)
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">{{ t('updateCleanup.title') }}</h1>
      <p class="content-subtitle">{{ t('updateCleanup.subtitle') }}</p>
    </div>
    <div class="content-body" style="flex:0 0 auto; padding-bottom:0;">
      <div class="stat-grid">
        <StatCard :label="t('status.updatable')" :value="totalOutdated" :sub="t('updateCleanup.hasNewVersion')" tone="accent" />
        <StatCard label="Formulae" :value="updateStore.formulae.length" :sub="t('updateCleanup.formula')" />
        <StatCard label="Casks" :value="updateStore.casks.length" :sub="t('updateCleanup.cask')" />
        <StatCard label="Pinned" :value="totalPinned" :sub="t('updateCleanup.pinned')" tone="warn" />
      </div>
    </div>
    <div class="toolbar">
      <SegmentedControl v-model="tab" :options="tabs" />
      <div class="toolbar-spacer">
        <BrewButton @click="refresh"><RefreshCw :size="14" />{{ t('common.refresh') }}</BrewButton>
        <BrewButton variant="primary" :disabled="totalOutdated === 0" @click="run('upgrade all', () => updateStore.upgradeAll())"><Upload :size="14" />{{ t('update.upgradeAll') }}</BrewButton>
      </div>
    </div>
    <div class="content-body">
      <p v-if="error" class="error-text">{{ error }}</p>
      <table v-if="tab === 'upgrades'" class="pkg-table">
        <thead><tr><th>{{ t('table.name') }}</th><th>{{ t('table.type') }}</th><th>{{ t('table.currentVersion') }}</th><th>{{ t('table.latestVersion') }}</th><th>{{ t('table.status') }}</th><th></th></tr></thead>
        <tbody>
          <tr v-for="item in updateStore.formulae" :key="`f-${item.name}`">
            <td class="pkg-name">{{ item.name }}</td><td>{{ t('installed.formulae') }}</td><td class="pkg-version">{{ item.installed_versions?.join(', ') }}</td><td class="pkg-version">{{ item.current_version }}</td><td><StatusPill :status="item.pinned ? 'pinned' : 'update'">{{ item.pinned ? t('status.pinned') : t('status.updatable') }}</StatusPill></td>
            <td style="text-align:right;">
              <div style="display:flex; gap:6px; justify-content:flex-end;">
                <BrewButton v-if="item.pinned" @click="run(`unpin ${item.name}`, () => BrewService.Unpin(item.name))">{{ t('actions.unpin') }}</BrewButton>
                <BrewButton v-else variant="primary" @click="run(`upgrade ${item.name}`, () => updateStore.upgrade(item.name))">{{ t('actions.update') }}</BrewButton>
              </div>
            </td>
          </tr>
          <tr v-for="item in updateStore.casks" :key="`c-${item.name}`">
            <td class="pkg-name">{{ item.name }}</td><td>{{ t('installed.casks') }}</td><td class="pkg-version">{{ item.installed_version }}</td><td class="pkg-version">{{ item.current_version }}</td><td><StatusPill status="update">{{ t('status.updatable') }}</StatusPill></td>
            <td style="text-align:right;"><BrewButton variant="primary" @click="run(`upgrade ${item.name}`, () => updateStore.upgrade(item.name))">{{ t('actions.update') }}</BrewButton></td>
          </tr>
        </tbody>
      </table>
      <div v-else-if="tab === 'cleanup'" style="display:flex; flex-direction:column; gap:12px;">
        <div class="toolbar" style="padding:0; border:0;"><BrewButton variant="danger" @click="run('brew cleanup', () => BrewService.Cleanup())"><Trash2 :size="14" />{{ t('cleanup.execute') }}</BrewButton></div>
        <TerminalPanel :lines="cleanupOutput ? cleanupOutput.split('\n') : []" :empty="t('updateCleanup.noCleanupCandidate')" />
      </div>
      <div v-else style="display:flex; flex-direction:column; gap:12px;">
        <div class="toolbar" style="padding:0; border:0;"><BrewButton variant="danger" @click="run('brew autoremove', () => BrewService.AutoRemove())"><Trash2 :size="14" />{{ t('updateCleanup.removeUnused') }}</BrewButton></div>
        <TerminalPanel :lines="autoremoveOutput ? autoremoveOutput.split('\n') : []" :empty="t('updateCleanup.noAutoremove')" />
      </div>
    </div>
  </section>
</template>
