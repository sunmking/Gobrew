<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
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
    error.value = err?.message || '加载包详情失败'
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
      <BrewButton variant="ghost" @click="router.push('/')"><ArrowLeft :size="14" />返回</BrewButton>
      <div style="min-width:0;">
        <h1 class="content-title pkg-name" style="font-size:22px;">{{ name }}</h1>
        <p class="content-subtitle">{{ info?.desc || (loading ? '加载中...' : '暂无描述') }}</p>
      </div>
      <div class="toolbar-spacer">
        <StatusPill :status="info?.pinned ? 'pinned' : info?.installed_version ? 'installed' : 'not-installed'">{{ info?.pinned ? '已锁定' : info?.installed_version ? '已安装' : '未安装' }}</StatusPill>
        <BrewButton v-if="!info?.installed_version" variant="primary" @click="runAction('install')"><Download :size="14" />安装</BrewButton>
        <BrewButton v-else-if="!info?.pinned" variant="primary" @click="runAction('upgrade')"><Upload :size="14" />更新</BrewButton>
        <BrewButton v-if="type === 'formula' && info?.installed_version && !info?.pinned" @click="runAction('pin')"><Lock :size="14" />锁定</BrewButton>
        <BrewButton v-if="type === 'formula' && info?.installed_version && info?.pinned" @click="runAction('unpin')"><Unlock :size="14" />解锁</BrewButton>
        <BrewButton v-if="info?.installed_version" variant="danger" @click="runAction('uninstall')"><Trash2 :size="14" />卸载</BrewButton>
      </div>
    </div>

    <div class="content-body">
      <p v-if="error" class="error-text">{{ error }}</p>
      <div v-else class="detail-grid">
        <div>
          <div class="section-title">依赖项</div>
          <div style="margin-bottom:16px;">
            <span v-for="dep in deps.length ? deps : info?.dependencies || []" :key="dep" class="dep-tag">{{ dep }}</span>
            <span v-if="!deps.length && !info?.dependencies?.length" class="content-subtitle">无</span>
          </div>
          <div class="section-title">被依赖于</div>
          <div style="margin-bottom:16px;">
            <span v-for="item in uses" :key="item" class="dep-tag">{{ item }}</span>
            <span v-if="!uses.length" class="content-subtitle">无</span>
          </div>
          <div class="section-title">命令输出</div>
          <TerminalPanel :lines="logStore.lines.slice(-80)" />
        </div>
        <aside style="display:flex; flex-direction:column; gap:12px;">
          <div class="detail-card">
            <div class="detail-card-title">版本信息</div>
            <div class="meta-row"><span class="meta-label">类型</span><span class="meta-value">{{ type }}</span></div>
            <div class="meta-row"><span class="meta-label">当前版本</span><span class="meta-value">{{ info?.installed_version || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">最新版本</span><span class="meta-value">{{ info?.current_version || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">锁定</span><span class="meta-value">{{ info?.pinned ? '是' : '否' }}</span></div>
            <div class="meta-row"><span class="meta-label">Tap</span><span class="meta-value">{{ info?.tap || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">License</span><span class="meta-value">{{ info?.license || '-' }}</span></div>
          </div>
          <div class="detail-card">
            <div class="detail-card-title">链接</div>
            <div class="meta-row"><span class="meta-label">主页</span><span class="meta-value">{{ info?.homepage || '-' }}</span></div>
            <div class="meta-row"><span class="meta-label">Token</span><span class="meta-value">{{ info?.token || info?.name || '-' }}</span></div>
          </div>
        </aside>
      </div>
    </div>
  </section>
</template>
