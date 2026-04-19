<script setup lang="ts">
defineProps<{
  emptyText?: string
}>()
</script>

<template>
  <div class="space-y-1.5">
    <slot />
    <div v-if="!$slots.default" class="pkg-empty py-8 text-center text-sm">{{ emptyText || 'No items' }}</div>
  </div>
</template>

<style scoped>
.pkg-empty {
  color: var(--color-text-tertiary);
}
</style>
