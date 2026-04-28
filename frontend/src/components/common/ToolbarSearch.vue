<script setup lang="ts">
import { Search } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()
const placeholderText = computed(() => props.placeholder || t('searchPalette.placeholder'))

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
      :placeholder="placeholderText"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    >
  </form>
</template>
