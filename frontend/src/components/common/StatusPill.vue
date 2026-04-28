<script setup lang="ts">
withDefaults(defineProps<{
  status?: 'installed' | 'update' | 'pinned' | 'not-installed' | 'running' | 'stopped' | 'error' | 'neutral'
}>(), {
  status: 'neutral',
})
</script>

<template>
  <span class="pill" :class="`pill-${status}`">
    <slot />
  </span>
</template>
