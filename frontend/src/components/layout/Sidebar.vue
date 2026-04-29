<script setup lang="ts">
import { computed, onMounted, ref, type Component } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { Archive, Boxes, FileText, RefreshCw, Server } from 'lucide-vue-next'
import { Capabilities } from '../../../bindings/changeme/services/brewservice'

const route = useRoute()
const { t } = useI18n()
const brewVersion = ref('')

const items = computed<Array<{ path: string; label: string; section: string; icon: Component; match: (path: string) => boolean }>>(() => [
  { path: '/', label: t('sidebarNav.allPackages'), section: t('sidebarNav.package'), icon: Boxes, match: path => path === '/' || path.startsWith('/packages/') },
  { path: '/update-cleanup', label: t('sidebarNav.updateCleanup'), section: t('sidebarNav.package'), icon: RefreshCw, match: path => path === '/update-cleanup' },
  { path: '/taps', label: t('sidebarNav.taps'), section: t('sidebarNav.resources'), icon: Archive, match: path => path === '/taps' },
  { path: '/services', label: t('sidebarNav.services'), section: t('sidebarNav.resources'), icon: Server, match: path => path === '/services' },
  { path: '/brewfile', label: t('sidebarNav.brewfile'), section: t('sidebarNav.resources'), icon: FileText, match: path => path === '/brewfile' },
])

const sections = computed(() => {
  const grouped = new Map<string, Array<{ path: string; label: string; section: string; icon: Component; match: (path: string) => boolean }>>()
  for (const item of items.value) {
    grouped.set(item.section, [...(grouped.get(item.section) ?? []), item])
  }
  return Array.from(grouped.entries())
})

onMounted(async () => {
  try {
    const caps = await Capabilities()
    const raw = String(caps?.brew_version || '').trim()
    brewVersion.value = raw.startsWith('Homebrew ') ? raw.slice('Homebrew '.length).trim() : raw
  } catch {
    brewVersion.value = ''
  }
})
</script>

<template>
  <nav class="app-sidebar">
    <template v-for="[section, sectionItems] in sections" :key="section">
      <div class="sidebar-section-label">{{ section }}</div>
      <router-link
        v-for="item in sectionItems"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: item.match(route.path) }"
      >
        <component :is="item.icon" :size="15" />
        <span>{{ item.label }}</span>
      </router-link>
    </template>
    <div class="sidebar-footer">
      <div class="brew-version">{{ brewVersion ? `Homebrew ${brewVersion}` : 'Homebrew' }}</div>
    </div>
  </nav>
</template>
