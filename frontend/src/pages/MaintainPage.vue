<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { BrushCleaning, BookOpen, PackageCheck } from 'lucide-vue-next'
import TapsPage from '@/pages/TapsPage.vue'
import CleanupPage from '@/pages/CleanupPage.vue'
import BundlePage from '@/pages/BundlePage.vue'

type MaintainPanel = 'taps' | 'cleanup' | 'bundle'
const LAST_PANEL_KEY = 'gobrew.maintain.last_panel'
const panelOrder: MaintainPanel[] = ['taps', 'cleanup', 'bundle']

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const tabs = computed(() => [
  {
    value: 'taps' as MaintainPanel,
    label: t('taps.title'),
    icon: BookOpen,
  },
  {
    value: 'cleanup' as MaintainPanel,
    label: t('cleanup.title'),
    icon: BrushCleaning,
  },
  {
    value: 'bundle' as MaintainPanel,
    label: t('bundle.title'),
    icon: PackageCheck,
  },
])

function normalizePanel(value: unknown): MaintainPanel {
  if (value === 'cleanup') return 'cleanup'
  if (value === 'bundle') return 'bundle'
  return 'taps'
}

const selectedPanel = computed<MaintainPanel>(() => {
  const raw = Array.isArray(route.query.panel) ? route.query.panel[0] : route.query.panel
  return normalizePanel(raw)
})

const tabRefs = ref<Array<HTMLButtonElement | null>>([])

function setTabRef(el: HTMLButtonElement | null, index: number) {
  tabRefs.value[index] = el
}

function readLastPanel(): MaintainPanel | null {
  if (typeof window === 'undefined') return null
  const stored = window.localStorage.getItem(LAST_PANEL_KEY)
  if (!stored) return null
  return normalizePanel(stored)
}

function writeLastPanel(panel: MaintainPanel) {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(LAST_PANEL_KEY, panel)
}

function switchPanel(nextPanel: MaintainPanel, replace = false) {
  if (nextPanel === selectedPanel.value) return
  const action = replace ? router.replace : router.push
  action({
    query: {
      ...route.query,
      panel: nextPanel,
    },
  })
}

function movePanel(offset: number): MaintainPanel {
  const currentIndex = panelOrder.indexOf(selectedPanel.value)
  const nextIndex = (currentIndex + offset + panelOrder.length) % panelOrder.length
  return panelOrder[nextIndex]
}

function onTabKeydown(event: KeyboardEvent) {
  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight' && event.key !== 'Home' && event.key !== 'End') return
  event.preventDefault()
  let nextPanel = selectedPanel.value
  if (event.key === 'ArrowLeft') nextPanel = movePanel(-1)
  if (event.key === 'ArrowRight') nextPanel = movePanel(1)
  if (event.key === 'Home') nextPanel = panelOrder[0]
  if (event.key === 'End') nextPanel = panelOrder[panelOrder.length - 1]
  switchPanel(nextPanel)
  nextTick(() => {
    const idx = panelOrder.indexOf(nextPanel)
    tabRefs.value[idx]?.focus()
  })
}

watch(
  () => route.query.panel,
  (panelParam) => {
    const raw = Array.isArray(panelParam) ? panelParam[0] : panelParam
    if (!raw) {
      const remembered = readLastPanel()
      if (remembered && remembered !== selectedPanel.value) {
        switchPanel(remembered, true)
      }
      return
    }
    const normalized = normalizePanel(raw)
    if (raw !== normalized) {
      switchPanel(normalized, true)
    }
  },
  { immediate: true },
)

watch(selectedPanel, (panel) => writeLastPanel(panel), { immediate: true })

const activePage = computed(() => {
  if (selectedPanel.value === 'cleanup') return CleanupPage
  if (selectedPanel.value === 'bundle') return BundlePage
  return TapsPage
})
</script>

<template>
  <div style="display:flex;flex-direction:column;gap:16px;height:100%;">
    <div class="content-header" style="margin-bottom:0;">
      <h1 class="content-title">{{ t('domains.maintain.title') }}</h1>
      <p class="content-subtitle">{{ t('domains.maintain.description') }}</p>
    </div>

    <div class="ui-segmented" role="tablist" :aria-label="t('domains.maintain.title')">
      <button
        v-for="(tab, index) in tabs"
        :key="tab.value"
        :ref="(el) => setTabRef(el as HTMLButtonElement | null, index)"
        class="ui-segmented-btn"
        :class="{ 'is-active': selectedPanel === tab.value }"
        role="tab"
        :aria-selected="selectedPanel === tab.value"
        :aria-controls="`maintain-panel-${tab.value}`"
        :tabindex="selectedPanel === tab.value ? 0 : -1"
        @click="switchPanel(tab.value)"
        @keydown="onTabKeydown"
      >
        <component :is="tab.icon" :size="14" />
        {{ tab.label }}
      </button>
    </div>

    <div :id="`maintain-panel-${selectedPanel}`" role="tabpanel" style="flex:1;min-height:0;overflow:hidden;">
      <KeepAlive>
        <component :is="activePage" embedded style="height:100%;" />
      </KeepAlive>
    </div>
  </div>
</template>
