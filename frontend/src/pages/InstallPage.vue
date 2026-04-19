<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { Compass, PackageOpen } from 'lucide-vue-next'
import ExplorePage from '@/pages/ExplorePage.vue'
import InstalledPage from '@/pages/InstalledPage.vue'

type InstallView = 'explore' | 'installed'
const tabOrder: InstallView[] = ['explore', 'installed']

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const tabs = computed(() => [
  { value: 'explore' as InstallView, label: t('explore.title'), icon: Compass },
  { value: 'installed' as InstallView, label: t('installed.title'), icon: PackageOpen },
])

function normalizeView(value: unknown): InstallView {
  return value === 'installed' ? 'installed' : 'explore'
}

const selectedView = computed<InstallView>({
  get() {
    const raw = Array.isArray(route.query.view) ? route.query.view[0] : route.query.view
    return normalizeView(raw)
  },
  set(nextView) {
    if (nextView === selectedView.value) return
    router.replace({
      query: {
        ...route.query,
        view: nextView,
      },
    })
  },
})

const tabRefs = ref<Array<HTMLButtonElement | null>>([])

function setTabRef(el: HTMLButtonElement | null, index: number) {
  tabRefs.value[index] = el
}

function switchView(nextView: InstallView) {
  selectedView.value = nextView
}

function moveView(offset: number): InstallView {
  const currentIndex = tabOrder.indexOf(selectedView.value)
  const nextIndex = (currentIndex + offset + tabOrder.length) % tabOrder.length
  return tabOrder[nextIndex]
}

function onTabKeydown(event: KeyboardEvent) {
  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight' && event.key !== 'Home' && event.key !== 'End') return
  event.preventDefault()
  let nextView = selectedView.value
  if (event.key === 'ArrowLeft') nextView = moveView(-1)
  if (event.key === 'ArrowRight') nextView = moveView(1)
  if (event.key === 'Home') nextView = tabOrder[0]
  if (event.key === 'End') nextView = tabOrder[tabOrder.length - 1]
  switchView(nextView)
  nextTick(() => {
    const idx = tabOrder.indexOf(nextView)
    tabRefs.value[idx]?.focus()
  })
}

const activePage = computed(() => (selectedView.value === 'installed' ? InstalledPage : ExplorePage))
</script>

<template>
  <div style="display:flex;flex-direction:column;gap:16px;height:100%;">
    <div class="content-header" style="margin-bottom:0;">
      <h1 class="content-title">{{ t('domains.install.title') }}</h1>
      <p class="content-subtitle">{{ t('domains.install.description') }}</p>
    </div>

    <div class="ui-segmented" role="tablist" :aria-label="t('install.title')">
      <button
        v-for="(tab, index) in tabs"
        :key="tab.value"
        :ref="(el) => setTabRef(el as HTMLButtonElement | null, index)"
        class="ui-segmented-btn"
        :class="{ 'is-active': selectedView === tab.value }"
        role="tab"
        :aria-selected="selectedView === tab.value"
        :aria-controls="`install-panel-${tab.value}`"
        :tabindex="selectedView === tab.value ? 0 : -1"
        @click="switchView(tab.value)"
        @keydown="onTabKeydown"
      >
        <component :is="tab.icon" :size="14" />
        {{ tab.label }}
      </button>
    </div>

    <div :id="`install-panel-${selectedView}`" role="tabpanel" style="flex:1;min-height:0;overflow:hidden;">
      <component :is="activePage" embedded style="height:100%;" />
    </div>
  </div>
</template>
