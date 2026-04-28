<script setup lang="ts">
defineProps<{
  modelValue: string
  options: Array<{ label: string; value: string }>
}>()

const emit = defineEmits<{ (event: 'update:modelValue', value: string): void }>()
</script>

<template>
  <div class="segmented-control">
    <button
      v-for="option in options"
      :key="option.value"
      class="segment"
      :class="{ active: option.value === modelValue }"
      type="button"
      @click="emit('update:modelValue', option.value)"
    >
      {{ option.label }}
    </button>
  </div>
</template>
