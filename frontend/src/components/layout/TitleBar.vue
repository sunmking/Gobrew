<script setup lang="ts">
import { Bell, Moon, Sun } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const theme = ref(document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light')
const isDark = computed(() => theme.value === 'dark')

function toggleTheme() {
  theme.value = isDark.value ? 'light' : 'dark'
  document.documentElement.dataset.theme = theme.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('gobrew-theme', theme.value)
}
</script>

<template>
  <header class="titlebar">
    <div class="traffic-lights" aria-hidden="true">
      <span class="dot dot-close" />
      <span class="dot dot-min" />
      <span class="dot dot-max" />
    </div>
    <div class="titlebar-title">Gobrew</div>
    <div class="titlebar-actions">
      <button class="btn-icon" type="button" :title="isDark ? '切换浅色' : '切换深色'" @click="toggleTheme">
        <Moon v-if="!isDark" :size="14" />
        <Sun v-else :size="14" />
      </button>
      <button class="btn-icon" type="button" title="通知">
        <Bell :size="14" />
      </button>
    </div>
  </header>
</template>
