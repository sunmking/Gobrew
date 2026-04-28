<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Plus, RefreshCw, Trash2 } from 'lucide-vue-next'
import ToolbarSearch from '@/components/common/ToolbarSearch.vue'
import BrewButton from '@/components/common/BrewButton.vue'
import StatCard from '@/components/common/StatCard.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import { useTapsStore } from '@/stores/taps'

const tapsStore = useTapsStore()
const query = ref('')
const newTap = ref('')
const selected = ref('')

const filtered = computed(() => tapsStore.taps.filter(tap => !query.value || tap.name.toLowerCase().includes(query.value.toLowerCase()) || tap.remote.toLowerCase().includes(query.value.toLowerCase())))
const detail = computed(() => tapsStore.selectedDetail)

async function select(name: string) {
  selected.value = name
  await tapsStore.details(name)
}

async function addTap() {
  const name = newTap.value.trim()
  if (!name) return
  await tapsStore.add(name)
  newTap.value = ''
}

onMounted(tapsStore.fetchTaps)
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">Taps</h1>
      <p class="content-subtitle">管理 Homebrew 第三方仓库</p>
    </div>
    <div class="toolbar">
      <ToolbarSearch v-model="query" placeholder="搜索 tap 名称或 URL..." />
      <input v-model="newTap" class="search-input" style="max-width:260px; padding-left:10px;" placeholder="homebrew/cask-fonts" @keyup.enter="addTap">
      <BrewButton variant="primary" @click="addTap"><Plus :size="14" />添加 Tap</BrewButton>
      <BrewButton @click="tapsStore.fetchTaps"><RefreshCw :size="14" />刷新</BrewButton>
    </div>
    <div class="content-body">
      <div class="stat-grid">
        <StatCard label="Taps" :value="tapsStore.taps.length" sub="已配置仓库" />
        <StatCard label="Formulae" :value="detail?.formula_names?.length || 0" sub="所选 tap" tone="accent" />
        <StatCard label="Casks" :value="detail?.cask_tokens?.length || 0" sub="所选 tap" />
        <StatCard label="Branch" :value="detail?.branch || '-'" sub="所选 tap" />
      </div>
      <div style="display:grid; grid-template-columns:minmax(0,1fr) 320px; gap:16px;">
        <div>
          <div v-for="tap in filtered" :key="tap.name" class="detail-card" style="margin-bottom:8px; cursor:pointer;" @click="select(tap.name)">
            <div style="display:flex; align-items:center; gap:12px;">
              <div style="flex:1; min-width:0;">
                <div class="pkg-name">{{ tap.name }}</div>
                <div class="pkg-desc">{{ tap.remote || 'remote unavailable' }}</div>
              </div>
              <StatusPill :status="tap.custom_remote ? 'neutral' : 'installed'">{{ tap.custom_remote ? 'Custom' : 'Official' }}</StatusPill>
              <BrewButton variant="ghost" @click.stop="tapsStore.remove(tap.name)"><Trash2 :size="14" /></BrewButton>
            </div>
          </div>
          <div v-if="filtered.length === 0" class="empty-state">没有匹配的 Tap</div>
        </div>
        <aside class="detail-card">
          <div class="detail-card-title">Tap 详情</div>
          <template v-if="detail">
            <div class="meta-row"><span class="meta-label">名称</span><span class="meta-value">{{ detail.name }}</span></div>
            <div class="meta-row"><span class="meta-label">Formulae</span><span class="meta-value">{{ detail.formula_names.length }}</span></div>
            <div class="meta-row"><span class="meta-label">Casks</span><span class="meta-value">{{ detail.cask_tokens.length }}</span></div>
            <div class="meta-row"><span class="meta-label">Last commit</span><span class="meta-value">{{ detail.last_commit || '-' }}</span></div>
          </template>
          <div v-else class="content-subtitle">选择一个 tap 查看详情</div>
        </aside>
      </div>
    </div>
  </section>
</template>
