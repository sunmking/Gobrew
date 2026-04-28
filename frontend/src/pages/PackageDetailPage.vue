<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Download, Lock, Trash2, Unlock, Upload } from 'lucide-vue-next'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import BrewButton from '@/components/common/BrewButton.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import TerminalPanel from '@/components/common/TerminalPanel.vue'
import { useLogStore } from '@/stores/log'
import type { PackageInfoResult, PackageType } from '@/types/brew'

const route = useRoute()
const router = useRouter()
const logStore = useLogStore()
const { t } = useI18n()
const info = ref<PackageInfoResult | null>(null)
const deps = ref<string[]>([])
const uses = ref<string[]>([])
const loading = ref(false)
const error = ref('')

const name = computed(() => String(route.params.name || ''))
const type = computed<PackageType>(() => route.params.type === 'cask' ? 'cask' : 'formula')

async function load() {
  loading.value = true
  error.value = ''
  deps.value = []
  uses.value = []
  try {
    info.value = (await BrewService.PackageInfo(name.value, type.value) as unknown) as PackageInfoResult
    if (type.value === 'formula') {
      deps.value = await BrewService.Deps(name.value, false)
      uses.value = await BrewService.Uses(name.value, true, false)
    }
  } catch (err: any) {
    error.value = err?.message || t('messages.operationFailed')
  } finally {
    loading.value = false
  }
}

async function runAction(action: 'install' | 'upgrade' | 'uninstall' | 'pin' | 'unpin') {
  logStore.startListening(`${action} ${name.value}`)
  try {
    if (action === 'install') await BrewService.Install(name.value)
    if (action === 'upgrade') await BrewService.Upgrade(name.value)
    if (action === 'pin') await BrewService.Pin(name.value)
    if (action === 'unpin') await BrewService.Unpin(name.value)
    if (action === 'uninstall') {
      if (type.value === 'cask') await BrewService.UninstallCask(name.value, false)
      else await BrewService.Uninstall(name.value)
    }
    await load()
  } finally {
    logStore.stopListening()
  }
}

watch(() => route.fullPath, load)
onMounted(load)
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <BrewButton variant="ghost" @click="router.push('/')"><ArrowLeft :size="14" />{{ t('actions.back') }}</BrewButton>
      <div style="min-width:0;">
        <h1 class="content-title pkg-name" style="font-size:22px;">{{ name }}</h1>
        <p class="content-subtitle">{{ info?.desc || (loading ? t('common.loading') : t('common.noDetails')) }}</p>
      </div>
      <div class="toolbar-spacer">
        <StatusPill :status="info?.pinned ? 'pinned' : info?.installed_version ? 'installed' : 'not-installed'">{{ info?.pinned ? t('status.pinned') : info?.installed_version ? t('common.installed') : t('status.notInstalled') }}</StatusPill>
        <BrewButton v-if="!info?.installed_version" variant="primary" @click="runAction('install')"><Download :size="14" />{{ t('actions.install') }}</BrewButton>
        <BrewButton v-else-if="!info?.pinned" variant="primary" @click="runAction('upgrade')"><Upload :size="14" />{{ t('actions.update') }}</BrewButton>
        <BrewButton v-if="type === 'formula' && info?.installed_version && !info?.pinned" @click="runAction('pin')"><Lock :size="14" />{{ t('actions.pin') }}</BrewButton>
        <BrewButton v-if="type === 'formula' && info?.installed_version && info?.pinned" @click="runAction('unpin')"><Unlock :size="14" />{{ t('actions.unpin') }}</BrewButton>
        <BrewButton v-if="info?.installed_version" variant="danger" @click="runAction('uninstall')"><Trash2 :size="14" />{{ t('actions.uninstall') }}</BrewButton>
      </div>
    </div>

    <div class="content-body">
      <p v-if="error" class="error-text">{{ error }}</p>
      <div v-else class="detail-grid">
        <div>
          <div class="section-title">{{ t('packageDetail.dependencies') }}</div>
          <div style="margin-bottom:16px;">
            <span v-for="dep in deps.length ? deps : info?.dependencies || []" :key="dep" class="dep-tag">{{ dep }}</span>
            <span v-if="!deps.length && !info?.dependencies?.length" class="content-subtitle">{{ t('packageDetail.none') }}</span>
          </div>
          <div class="section-title">{{ t('packageDetail.requiredBy') }}</div>
          <div style="margin-bottom:16px;">
            <span v-for="item in uses" :key="item" class="dep-tag">{{ item }}</span>
            <span v-if="!uses.length" class="content-subtitle">{{ t('packageDetail.none') }}</span>
          </div>
          <div class="section-title">{{ t('bundlePage.commandOutput') }}</div>
          <TerminalPanel :lines="logStore.lines.slice(-80)" />
        </div>
        <aside style="display:flex; flex-direction:column; gap:12px;">
          <div class="detail-card">
            <div class="detail-card-title">{{ t('packageDetail.versionInfo') }}</div>
            <div class="meta-row"><span class="meta-label">{{ t('table.type') }}</span><span class="meta-value">{{ type }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('table.currentVersion') }}</span><span class="meta-value">{{ info?.installed_version || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('table.latestVersion') }}</span><span class="meta-value">{{ info?.current_version || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('status.pinned') }}</span><span class="meta-value">{{ info?.pinned ? t('common.yes') : t('common.no') }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('packageDetail.tap') }}</span><span class="meta-value">{{ info?.tap || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('packageDetail.license') }}</span><span class="meta-value">{{ info?.license || '-' }}</span></div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">{{ t('packageDetail.links') }}</div>
            <div class="meta-row"><span class="meta-label">{{ t('packageDetail.homepage') }}</span><span class="meta-value">{{ info?.homepage || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('packageDetail.token') }}</span><span class="meta-value">{{ info?.token || info?.name || '-' }}</span></div>
          </div>
        </aside>
      </div>
    </div>
  </section>
</template>
