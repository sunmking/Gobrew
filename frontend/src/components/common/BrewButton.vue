<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  disabled?: boolean
  type?: 'button' | 'submit'
}>(), {
  variant: 'secondary',
  disabled: false,
  type: 'button',
})
</script>

<template>
  <button class="btn" :class="`btn-${variant}`" :type="type" :disabled="disabled">
    <slot />
  </button>
</template>
