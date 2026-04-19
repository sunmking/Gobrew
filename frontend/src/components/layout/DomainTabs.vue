<script setup lang="ts">
import { useI18n } from 'vue-i18n'

export type DomainTabOption = {
  value: string
  label: string
  disabled?: boolean
}

const props = defineProps<{
  modelValue: string
  tabs: DomainTabOption[]
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
}>()
const { t } = useI18n()

function selectTab(tab: DomainTabOption) {
  if (tab.disabled || tab.value === props.modelValue) return
  emit('update:modelValue', tab.value)
}
</script>

<template>
  <div class="flex items-center gap-4 border-b border-[var(--color-border)]" role="tablist" :aria-label="t('common.domainViews')">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      type="button"
      role="tab"
      :aria-selected="modelValue === tab.value"
      :disabled="tab.disabled"
      class="relative pb-2 pt-1 text-sm font-medium transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-accent)] focus-visible:ring-offset-2"
      :class="[
        modelValue === tab.value
          ? 'text-[var(--color-accent)]'
          : 'text-[var(--color-text-secondary)] hover:text-[var(--color-text-primary)]',
        tab.disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer',
      ]"
      @click="selectTab(tab)"
    >
      {{ tab.label }}
      <span
        v-if="modelValue === tab.value"
        class="absolute bottom-0 left-0 right-0 h-[2px] rounded-full bg-[var(--color-accent)]"
      />
    </button>
  </div>
</template>
