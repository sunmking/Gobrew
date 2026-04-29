<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, RefreshCw, Trash2 } from 'lucide-vue-next'
import ToolbarSearch from '@/components/common/ToolbarSearch.vue'
import BrewButton from '@/components/common/BrewButton.vue'
import StatCard from '@/components/common/StatCard.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import { useTapsStore } from '@/stores/taps'

const tapsStore = useTapsStore()
const { t } = useI18n()
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
      <h1 class="content-title">{{ t('taps.title') }}</h1>
      <p class="content-subtitle">{{ t('tapsPage.subtitle') }}</p>
    </div>
    <div class="toolbar">
      <ToolbarSearch v-model="query" :placeholder="t('tapsPage.searchPlaceholder')" />
      <input v-model="newTap" class="search-input" style="max-width:260px; padding-left:10px;" :placeholder="t('tapsPage.newTapPlaceholder')" @keyup.enter="addTap">
      <BrewButton variant="primary" @click="addTap"><Plus :size="14" />{{ t('taps.add') }}</BrewButton>
      <BrewButton @click="tapsStore.fetchTaps"><RefreshCw :size="14" />{{ t('common.refresh') }}</BrewButton>
    </div>
    <div class="content-body">
      <div class="stat-grid">
        <StatCard :label="t('tapsPage.stats.taps')" :value="tapsStore.taps.length" :sub="t('tapsPage.configuredRepo')" />
        <StatCard :label="t('tapsPage.stats.formulae')" :value="detail?.formula_names?.length || 0" :sub="t('tapsPage.selectedTap')" tone="accent" />
        <StatCard :label="t('tapsPage.stats.casks')" :value="detail?.cask_tokens?.length || 0" :sub="t('tapsPage.selectedTap')" />
        <StatCard :label="t('tapsPage.stats.branch')" :value="detail?.branch || '-'" :sub="t('tapsPage.selectedTap')" />
      </div>
      <div style="display:grid; grid-template-columns:minmax(0,1fr) 320px; gap:16px;">
        <div>
          <div v-for="tap in filtered" :key="tap.name" class="detail-card" style="margin-bottom:8px; cursor:pointer;" @click="select(tap.name)">
            <div style="display:flex; align-items:center; gap:12px;">
              <div style="flex:1; min-width:0;">
                <div class="pkg-name">{{ tap.name }}</div>
                <div class="pkg-desc">{{ tap.remote || t('tapsPage.remoteUnavailable') }}</div>
              </div>
              <StatusPill :status="tap.custom_remote ? 'neutral' : 'installed'">{{ tap.custom_remote ? t('taps.detail.customRemote') : t('tapsPage.official') }}</StatusPill>
              <BrewButton variant="ghost" @click.stop="tapsStore.remove(tap.name)"><Trash2 :size="14" /></BrewButton>
            </div>
          </div>
          <div v-if="filtered.length === 0" class="empty-state">{{ t('tapsPage.noMatch') }}</div>
        </div>
        <aside class="detail-card">
          <div class="detail-card-title">{{ t('tapsPage.detail') }}</div>
          <template v-if="detail">
            <div class="meta-row"><span class="meta-label">{{ t('table.name') }}</span><span class="meta-value">{{ detail.name }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('tapsPage.stats.formulae') }}</span><span class="meta-value">{{ detail.formula_names.length }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('tapsPage.stats.casks') }}</span><span class="meta-value">{{ detail.cask_tokens.length }}</span></div>
            <div class="meta-row"><span class="meta-label">{{ t('taps.detail.lastCommit') }}</span><span class="meta-value">{{ detail.last_commit || '-' }}</span></div>
          </template>
          <div v-else class="content-subtitle">{{ t('tapsPage.pickOne') }}</div>
        </aside>
      </div>
    </div>
  </section>
</template>
