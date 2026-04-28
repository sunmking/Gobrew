<script setup lang="ts">
withDefaults(defineProps<{
  label: string
  value: string | number
  sub?: string
  tone?: 'default' | 'accent' | 'warn' | 'success' | 'danger'
}>(), {
  sub: '',
  tone: 'default',
})
</script>

<template>
  <div class="stat-card">
    <div class="stat-label">{{ label }}</div>
    <div class="stat-value" :class="tone">{{ value }}</div>
    <div v-if="sub" class="stat-sub">{{ sub }}</div>
  </div>
</template>
