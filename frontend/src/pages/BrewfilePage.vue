<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Clipboard, Download, FileDown, RefreshCw, Save, Trash2 } from 'lucide-vue-next'
import BrewButton from '@/components/common/BrewButton.vue'
import StatCard from '@/components/common/StatCard.vue'
import TerminalPanel from '@/components/common/TerminalPanel.vue'
import { useBundleStore } from '@/stores/bundle'
import { useLogStore } from '@/stores/log'
import { useSettingsStore } from '@/stores/settings'

const bundleStore = useBundleStore()
const logStore = useLogStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()
const path = ref('')
const content = ref('')

const missing = computed(() => bundleStore.checkResult?.missing ?? [])
const stats = computed(() => {
  const lines = content.value.split('\n').map(line => line.trim())
  return {
    taps: lines.filter(line => line.startsWith('tap ')).length,
    brews: lines.filter(line => line.startsWith('brew ')).length,
    casks: lines.filter(line => line.startsWith('cask ')).length,
  }
})

async function read() {
  await bundleStore.readBrewfile(path.value || undefined)
  content.value = bundleStore.brewfileContent
}

async function run(label: string, fn: () => Promise<void>) {
  logStore.startListening(label)
  try {
    await fn()
  } finally {
    logStore.stopListening()
  }
}

async function save() {
  await bundleStore.write(path.value || undefined, content.value)
}

async function dump() {
  await run('brew bundle dump', async () => {
    await bundleStore.dump(path.value || undefined, true)
    content.value = bundleStore.brewfileContent
  })
}

onMounted(async () => {
  if (!settingsStore.initialized) {
    await settingsStore.load()
  }
  path.value = settingsStore.brewFilePath
  await read()
})

async function updatePath(value: string) {
  path.value = value
  await settingsStore.setBrewFilePath(value)
}
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">{{ t('bundlePage.title') }}</h1>
      <p class="content-subtitle">{{ t('bundlePage.subtitle') }}</p>
    </div>
    <div class="content-body">
      <div class="stat-grid">
        <StatCard :label="t('bundlePage.taps')" :value="stats.taps" />
        <StatCard :label="t('bundlePage.brews')" :value="stats.brews" />
        <StatCard :label="t('bundlePage.casks')" :value="stats.casks" />
        <StatCard :label="t('bundlePage.pendingInstall')" :value="missing.length" tone="accent" />
      </div>
      <div style="display:grid; grid-template-columns:minmax(0,1fr) 320px; gap:20px;">
        <div>
          <div class="toolbar" style="padding:0 0 10px; border:0;">
            <input :value="path" class="search-input" style="max-width:360px; padding-left:10px;" :placeholder="t('bundlePage.defaultBrewfilePath')" @change="updatePath(($event.target as HTMLInputElement).value)">
            <BrewButton @click="read"><RefreshCw :size="14" />{{ t('common.refresh') }}</BrewButton>
            <BrewButton variant="primary" @click="save"><Save :size="14" />{{ t('bundlePage.save') }}</BrewButton>
          </div>
          <textarea v-model="content" spellcheck="false" style="width:100%; min-height:420px; resize:vertical; background:var(--surface); color:var(--fg); border:1px solid var(--border); border-radius:var(--radius); padding:12px 14px; font-family:var(--font-mono); font-size:12.5px; line-height:1.7; outline:none;" />
        </div>
        <aside style="display:flex; flex-direction:column; gap:12px;">
          <div class="detail-card">
            <div class="detail-card-title">{{ t('bundlePage.actions') }}</div>
            <div style="display:flex; flex-direction:column; gap:8px;">
              <BrewButton variant="primary" @click="run('brew bundle install', () => bundleStore.restore(path || undefined))"><Download :size="14" />{{ t('bundlePage.restore') }}</BrewButton>
              <BrewButton @click="dump"><FileDown :size="14" />{{ t('bundlePage.exportFromEnv') }}</BrewButton>
              <BrewButton @click="bundleStore.check(path || undefined)"><Clipboard :size="14" />{{ t('bundlePage.checkDiff') }}</BrewButton>
              <BrewButton variant="danger" @click="run('brew bundle cleanup', () => bundleStore.cleanup(path || undefined, true))"><Trash2 :size="14" />{{ t('bundlePage.cleanupExtra') }}</BrewButton>
            </div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">{{ t('bundlePage.installDiffPreview') }}</div>
            <div v-if="missing.length">
              <div v-for="item in missing" :key="item" style="font-family:var(--font-mono); font-size:12px; color:var(--success);">+ {{ item }}</div>
            </div>
            <div v-else class="content-subtitle">{{ t('bundlePage.noMissingDeps') }}</div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">{{ t('bundlePage.commandOutput') }}</div>
            <TerminalPanel :lines="bundleStore.cleanupPreviewResult?.output ? bundleStore.cleanupPreviewResult.output.split('\n') : logStore.lines.slice(-40)" />
          </div>
        </aside>
      </div>
    </div>
  </section>
</template>
