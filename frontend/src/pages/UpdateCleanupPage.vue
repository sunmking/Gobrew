<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
const tab = ref('upgrades')
const cleanupOutput = ref('')
const autoremoveOutput = ref('')
const error = ref('')

const totalOutdated = computed(() => updateStore.formulae.length + updateStore.casks.length)
const tabs = [
  { label: '待更新', value: 'upgrades' },
  { label: '清理', value: 'cleanup' },
  { label: '依赖缓存', value: 'cache' },
]

async function refresh() {
  error.value = ''
  try {
    await updateStore.forceRefresh()
    const [cleanup, autoremove] = await Promise.all([BrewService.CleanupPreview(), BrewService.AutoRemovePreview()])
    cleanupOutput.value = cleanup?.output || ''
    autoremoveOutput.value = autoremove?.output || ''
  } catch (err: any) {
    error.value = err?.message || '刷新失败'
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
      <h1 class="content-title">更新 & 清理</h1>
      <p class="content-subtitle">批量管理待更新包，清理旧版本和不再需要的依赖</p>
    </div>
    <div class="content-body" style="flex:0 0 auto; padding-bottom:0;">
      <div class="stat-grid">
        <StatCard label="可更新" :value="totalOutdated" sub="个包有新版本" tone="accent" />
        <StatCard label="Formulae" :value="updateStore.formulae.length" sub="命令行包" />
        <StatCard label="Casks" :value="updateStore.casks.length" sub="应用包" />
        <StatCard label="清理候选" :value="cleanupOutput ? cleanupOutput.split('\n').filter(Boolean).length : 0" sub="dry-run 输出行" tone="warn" />
      </div>
    </div>
    <div class="toolbar">
      <SegmentedControl v-model="tab" :options="tabs" />
      <div class="toolbar-spacer">
        <BrewButton @click="refresh"><RefreshCw :size="14" />刷新</BrewButton>
        <BrewButton variant="primary" :disabled="totalOutdated === 0" @click="run('upgrade all', () => updateStore.upgradeAll())"><Upload :size="14" />全部更新</BrewButton>
      </div>
    </div>
    <div class="content-body">
      <p v-if="error" class="error-text">{{ error }}</p>
      <table v-if="tab === 'upgrades'" class="pkg-table">
        <thead><tr><th>包名</th><th>类型</th><th>当前版本</th><th>最新版本</th><th>状态</th><th></th></tr></thead>
        <tbody>
          <tr v-for="item in updateStore.formulae" :key="`f-${item.name}`">
            <td class="pkg-name">{{ item.name }}</td><td>Formula</td><td class="pkg-version">{{ item.installed_versions?.join(', ') }}</td><td class="pkg-version">{{ item.current_version }}</td><td><StatusPill status="update">可更新</StatusPill></td>
            <td style="text-align:right;"><BrewButton variant="primary" @click="run(`upgrade ${item.name}`, () => updateStore.upgrade(item.name))">更新</BrewButton></td>
          </tr>
          <tr v-for="item in updateStore.casks" :key="`c-${item.name}`">
            <td class="pkg-name">{{ item.name }}</td><td>Cask</td><td class="pkg-version">{{ item.installed_version }}</td><td class="pkg-version">{{ item.current_version }}</td><td><StatusPill status="update">可更新</StatusPill></td>
            <td style="text-align:right;"><BrewButton variant="primary" @click="run(`upgrade ${item.name}`, () => updateStore.upgrade(item.name))">更新</BrewButton></td>
          </tr>
        </tbody>
      </table>
      <div v-else-if="tab === 'cleanup'" style="display:flex; flex-direction:column; gap:12px;">
        <div class="toolbar" style="padding:0; border:0;"><BrewButton variant="danger" @click="run('brew cleanup', () => BrewService.Cleanup())"><Trash2 :size="14" />执行清理</BrewButton></div>
        <TerminalPanel :lines="cleanupOutput ? cleanupOutput.split('\n') : []" empty="暂无清理候选" />
      </div>
      <div v-else style="display:flex; flex-direction:column; gap:12px;">
        <div class="toolbar" style="padding:0; border:0;"><BrewButton variant="danger" @click="run('brew autoremove', () => BrewService.AutoRemove())"><Trash2 :size="14" />移除无用依赖</BrewButton></div>
        <TerminalPanel :lines="autoremoveOutput ? autoremoveOutput.split('\n') : []" empty="暂无可移除依赖" />
      </div>
    </div>
  </section>
</template>
