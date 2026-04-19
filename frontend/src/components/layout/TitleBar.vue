<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    description: string
  }>(),
  {
    title: 'Gobrew',
    description: '',
  },
)
</script>

<template>
  <header class="shell-topbar">
    <div>
      <h1 class="shell-title">{{ title }}</h1>
      <p v-if="description" class="shell-description">{{ description }}</p>
    </div>
    <div class="shell-brand">
      <span class="shell-brand-mark" aria-hidden="true">🍺</span>
      <span class="shell-brand-wordmark">Gobrew</span>
    </div>
  </header>
</template>
