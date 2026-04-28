<script setup lang="ts">
import { Search } from 'lucide-vue-next'

withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
}>(), {
  placeholder: '搜索...',
})

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
  (event: 'submit'): void
}>()
</script>

<template>
  <form class="search-wrap" @submit.prevent="emit('submit')">
    <Search class="search-icon" :size="14" />
    <input
      class="search-input"
      :value="modelValue"
      :placeholder="placeholder"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    >
  </form>
</template>
