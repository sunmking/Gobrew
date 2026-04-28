<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Clipboard, Download, FileDown, RefreshCw, Save, Trash2 } from 'lucide-vue-next'
import BrewButton from '@/components/common/BrewButton.vue'
import StatCard from '@/components/common/StatCard.vue'
import TerminalPanel from '@/components/common/TerminalPanel.vue'
import { useBundleStore } from '@/stores/bundle'
import { useLogStore } from '@/stores/log'

const bundleStore = useBundleStore()
const logStore = useLogStore()
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
  await bundleStore.readBrewfile(path.value)
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
  await bundleStore.write(path.value, content.value)
}

async function dump() {
  await run('brew bundle dump', async () => {
    await bundleStore.dump(path.value, true)
    content.value = bundleStore.brewfileContent
  })
}

onMounted(read)
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">Brewfile</h1>
      <p class="content-subtitle">从 Brewfile 恢复环境，或将当前安装导出为 Brewfile</p>
    </div>
    <div class="content-body">
      <div class="stat-grid">
        <StatCard label="Taps" :value="stats.taps" />
        <StatCard label="Brews" :value="stats.brews" />
        <StatCard label="Casks" :value="stats.casks" />
        <StatCard label="待安装" :value="missing.length" tone="accent" />
      </div>
      <div style="display:grid; grid-template-columns:minmax(0,1fr) 320px; gap:20px;">
        <div>
          <div class="toolbar" style="padding:0 0 10px; border:0;">
            <input v-model="path" class="search-input" style="max-width:360px; padding-left:10px;" placeholder="默认 ~/Brewfile">
            <BrewButton @click="read"><RefreshCw :size="14" />读取</BrewButton>
            <BrewButton variant="primary" @click="save"><Save :size="14" />保存</BrewButton>
          </div>
          <textarea v-model="content" spellcheck="false" style="width:100%; min-height:420px; resize:vertical; background:var(--surface); color:var(--fg); border:1px solid var(--border); border-radius:var(--radius); padding:12px 14px; font-family:var(--font-mono); font-size:12.5px; line-height:1.7; outline:none;" />
        </div>
        <aside style="display:flex; flex-direction:column; gap:12px;">
          <div class="detail-card">
            <div class="detail-card-title">操作</div>
            <div style="display:flex; flex-direction:column; gap:8px;">
              <BrewButton variant="primary" @click="run('brew bundle install', () => bundleStore.restore(path))"><Download :size="14" />brew bundle install</BrewButton>
              <BrewButton @click="dump"><FileDown :size="14" />从环境导出</BrewButton>
              <BrewButton @click="bundleStore.check(path)"><Clipboard :size="14" />检查差异</BrewButton>
              <BrewButton variant="danger" @click="run('brew bundle cleanup', () => bundleStore.cleanup(path, true))"><Trash2 :size="14" />清理多余包</BrewButton>
            </div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">安装差异预览</div>
            <div v-if="missing.length">
              <div v-for="item in missing" :key="item" style="font-family:var(--font-mono); font-size:12px; color:var(--success);">+ {{ item }}</div>
            </div>
            <div v-else class="content-subtitle">暂无缺失依赖</div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">命令输出</div>
            <TerminalPanel :lines="bundleStore.cleanupPreviewResult?.output ? bundleStore.cleanupPreviewResult.output.split('\n') : logStore.lines.slice(-40)" />
          </div>
        </aside>
      </div>
    </div>
  </section>
</template>
