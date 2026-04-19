<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Trash2, Download, Shield } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import { useInstalledStore } from '@/stores/installed'
import { useUpdateStore } from '@/stores/update'
import { useLogStore } from '@/stores/log'
import Toast from '@/components/common/Toast.vue'
import PackageIcon from '@/components/icons/PackageIcon.vue'

const { t } = useI18n()
const installedStore = useInstalledStore()
const updateStore = useUpdateStore()
const logStore = useLogStore()
const toastRef = ref<any>(null)

const totalInstalled = computed(() => installedStore.formulae.length + installedStore.casks.length)
const totalOutdated = computed(() => updateStore.formulae.length + updateStore.casks.length)
const doctorStatus = ref<'idle' | 'ok' | 'issues'>('idle')

async function refreshAll() {
  await Promise.all([installedStore.fetchInstalled(), updateStore.fetchOutdated()])
}

async function upgradeAll() {
  logStore.startListening()
  try {
    await BrewService.UpgradeAll()
    toastRef.value?.show('success', t('home.allUpgraded'))
    await refreshAll()
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.upgradeFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function runCleanup() {
  logStore.startListening()
  try {
    await BrewService.Cleanup()
    toastRef.value?.show('success', t('home.cleanupDone'))
  } catch (error: any) {
    toastRef.value?.show('error', error?.message || t('messages.cleanupFailed'))
  } finally {
    logStore.stopListening()
  }
}

async function runDoctor() {
  logStore.startListening()
  try {
    await BrewService.Doctor()
    doctorStatus.value = 'ok'
    toastRef.value?.show('success', t('home.doctorOk'))
  } catch (error: any) {
    doctorStatus.value = 'issues'
    toastRef.value?.show('error', error?.message || t('messages.doctorFoundIssues'))
  } finally {
    logStore.stopListening()
  }
}

onMounted(refreshAll)
</script>

<template>
  <div class="home-page">
    <div class="content-header">
      <h1 class="content-title">{{ t('home.title') }}</h1>
      <p class="content-subtitle">{{ t('home.quickStatusDesc') }}</p>
    </div>

    <div class="content-section">
      <div class="stat-grid">
        <div class="stat-card">
          <div class="stat-value">{{ totalInstalled }}</div>
          <div class="stat-label">{{ t('home.installed') }}</div>
        </div>
        <div class="stat-card" style="border-left:3px solid var(--color-warning);">
          <div class="stat-value" style="color:var(--color-text);">{{ totalOutdated }}</div>
          <div class="stat-label">{{ t('home.outdated') }}</div>
        </div>
        <div class="stat-card" :style="doctorStatus === 'ok' ? 'border-left:3px solid var(--color-success)' : doctorStatus === 'issues' ? 'border-left:3px solid var(--color-danger)' : 'border-left:3px solid var(--color-border)'">
          <div class="stat-value" style="font-size:15px;">
            {{ doctorStatus === 'ok' ? '✓' : doctorStatus === 'issues' ? '✗' : '…' }}
          </div>
          <div class="stat-label">{{ doctorStatus === 'ok' ? t('home.doctorHealthy') : doctorStatus === 'issues' ? t('home.doctorIssues') : t('home.quickStatus') }}</div>
        </div>
      </div>
    </div>

    <div v-if="totalOutdated > 0" class="content-section">
      <h2 class="section-title">{{ t('home.outdated') }} ({{ totalOutdated }})</h2>
      <div class="card-group">
        <div v-for="item in updateStore.formulae" :key="item.name" class="card-row">
          <span style="display:flex;align-items:center;gap:8px;min-width:0;flex:1;">
            <PackageIcon :name="item.name" :size="28" />
            <div style="min-width:0;">
              <div class="card-row-title">{{ item.name }}</div>
              <div class="card-row-subtitle">{{ item.installed_versions?.[0] }}</div>
            </div>
          </span>
          <div style="display:flex;align-items:center;gap:10px;">
            <span class="card-row-detail">{{ item.current_version }}</span>
            <button class="btn-primary" @click="upgradeAll">{{ t('home.updateAll') }}</button>
          </div>
        </div>
        <div v-for="item in updateStore.casks" :key="item.name" class="card-row">
          <span style="display:flex;align-items:center;gap:8px;min-width:0;flex:1;">
            <PackageIcon :name="item.name" :size="28" />
            <div style="min-width:0;">
              <div class="card-row-title">{{ item.name }}</div>
              <div class="card-row-subtitle">{{ item.installed_version }}</div>
            </div>
          </span>
          <div style="display:flex;align-items:center;gap:10px;">
            <span class="card-row-detail">{{ item.current_version }}</span>
            <button class="btn-primary" @click="upgradeAll">{{ t('home.updateAll') }}</button>
          </div>
        </div>
      </div>
    </div>

    <div class="content-section">
      <div class="action-bar">
        <button class="btn-primary" :disabled="totalOutdated === 0" @click="upgradeAll">
          <Download :size="14" /> {{ t('home.updateAll') }}
        </button>
        <button class="btn-secondary" @click="runCleanup">
          <Trash2 :size="14" /> {{ t('home.cleanup') }}
        </button>
        <button class="btn-secondary" @click="refreshAll">
          <RefreshCw :size="14" /> {{ t('common.refresh') }}
        </button>
        <button class="btn-secondary" @click="runDoctor">
          <Shield :size="14" /> {{ t('home.runDoctor') }}
        </button>
      </div>
    </div>

    <Toast ref="toastRef" />
  </div>
</template>

<style scoped>
.home-page {
  height: 100%;
  overflow-y: auto;
}

.btn-primary {
  background: var(--color-accent);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--motion-fast) var(--ease);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.btn-primary:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--color-sidebar-hover);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 6px 14px;
  font-size: 13px;
  cursor: pointer;
  transition: background var(--motion-fast) var(--ease);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.btn-secondary:hover {
  background: var(--color-sidebar-active);
}
</style>
