<script setup lang="ts">
import { computed, type Component } from 'vue'
import { useRoute } from 'vue-router'
import { Archive, Boxes, FileText, RefreshCw, Server } from 'lucide-vue-next'

const route = useRoute()

const items: Array<{ path: string; label: string; section: string; icon: Component; match: (path: string) => boolean }> = [
  { path: '/', label: '所有包', section: '包管理', icon: Boxes, match: path => path === '/' || path.startsWith('/packages/') },
  { path: '/update-cleanup', label: '更新 & 清理', section: '包管理', icon: RefreshCw, match: path => path === '/update-cleanup' },
  { path: '/taps', label: 'Taps', section: '资源', icon: Archive, match: path => path === '/taps' },
  { path: '/services', label: 'Services', section: '资源', icon: Server, match: path => path === '/services' },
  { path: '/brewfile', label: 'Brewfile', section: '资源', icon: FileText, match: path => path === '/brewfile' },
]

const sections = computed(() => {
  const grouped = new Map<string, typeof items>()
  for (const item of items) {
    grouped.set(item.section, [...(grouped.get(item.section) ?? []), item])
  }
  return Array.from(grouped.entries())
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
      <div class="brew-version">Homebrew</div>
    </div>
  </nav>
</template>
