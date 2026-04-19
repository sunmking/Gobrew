<script setup lang="ts">
import { computed, ref, type Component } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { Home, Download, ArrowUpCircle, Wrench, Server, Settings } from 'lucide-vue-next'
import { DOMAIN_DISPLAY, DOMAIN_ENTRY_PATHS, DOMAIN_ORDER, resolveDomainFromPath, type Domain } from '@/config/domains'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const searchQuery = ref('')

const mainNav: Array<{ domain: Domain; icon: Component; path: string }> = [
  { domain: 'dashboard', icon: Home, path: '/' },
  { domain: 'install', icon: Download, path: '/install' },
  { domain: 'upgrade', icon: ArrowUpCircle, path: '/upgrade' },
  { domain: 'maintain', icon: Wrench, path: '/maintain' },
  { domain: 'services', icon: Server, path: '/services' },
]

const mainItems = computed(() => mainNav.map(item => ({
  ...item,
  label: t(DOMAIN_DISPLAY[item.domain].labelKey),
})))

const activeDomain = computed<Domain>(() => resolveDomainFromPath(route.path))
const activePath = computed(() => route.path)

function isActive(domain: Domain): boolean {
  return activeDomain.value === domain
}

const filteredMain = computed(() => {
  if (!searchQuery.value) return mainItems.value
  const q = searchQuery.value.toLowerCase()
  return mainItems.value.filter(i => i.label.toLowerCase().includes(q))
})
</script>

<template>
  <nav class="app-sidebar">
    <ul class="sidebar-nav">
      <li v-for="item in filteredMain" :key="item.domain">
        <router-link
          :to="item.path"
          class="sidebar-item"
          :class="{ 'is-active': isActive(item.domain) }"
        >
          <span class="sidebar-item-icon"><component :is="item.icon" :size="15" /></span>
          <span class="sidebar-item-label">{{ item.label }}</span>
        </router-link>
      </li>
    </ul>
    <div style="flex:1" />
    <div class="sidebar-divider" />
    <ul class="sidebar-nav">
      <li>
        <router-link
          to="/settings"
          class="sidebar-item"
          :class="{ 'is-active': activePath === '/settings' }"
        >
          <span class="sidebar-item-icon"><Settings :size="15" /></span>
          <span class="sidebar-item-label">{{ t('settings.title') }}</span>
        </router-link>
      </li>
    </ul>
  </nav>
</template>
