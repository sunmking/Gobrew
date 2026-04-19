<script setup lang="ts">
import { useI18n } from 'vue-i18n'

defineProps<{
  selectedCount: number
  allSelected: boolean
}>()

const emit = defineEmits<{
  selectAll: []
  clear: []
}>()

const { t } = useI18n()
</script>

<template>
  <div class="bulk-bar mb-3 flex flex-wrap items-center justify-between gap-2 rounded-lg p-2">
    <div class="bulk-text text-xs">
      {{ t('bulk.selected') }}: <span class="font-semibold">{{ selectedCount }}</span>
    </div>
    <div class="bulk-slot-actions flex items-center gap-2">
      <button
        class="bulk-toggle rounded px-2 py-1 text-xs"
        @click="allSelected ? emit('clear') : emit('selectAll')"
      >
        {{ allSelected ? t('bulk.clear') : t('bulk.selectAll') }}
      </button>
      <slot />
    </div>
  </div>
</template>

<style scoped>
.bulk-bar {
  background: var(--color-accent-light);
  border: 1px solid var(--color-accent);
}

.bulk-text {
  color: var(--color-accent);
}

.bulk-toggle {
  background: var(--color-group-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.bulk-toggle:hover {
  opacity: 0.85;
}

.bulk-slot-actions :deep(.btn-primary),
.bulk-slot-actions :deep(.btn-neutral),
.bulk-slot-actions :deep(.btn-danger) {
  border-radius: 0.25rem;
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  line-height: 1rem;
  transition: background-color 0.15s ease;
}

.bulk-slot-actions :deep(.btn-primary) {
  background: var(--color-accent);
  color: white;
  font-weight: 500;
}

.bulk-slot-actions :deep(.btn-primary:hover) {
  filter: brightness(1.08);
}

.bulk-slot-actions :deep(.btn-neutral) {
  background: var(--color-group-bg);
  color: var(--color-text-secondary);
}

.bulk-slot-actions :deep(.btn-neutral:hover) {
  opacity: 0.85;
}

.bulk-slot-actions :deep(.btn-danger) {
  background: var(--color-danger);
  color: white;
}

.bulk-slot-actions :deep(.btn-danger:hover) {
  filter: brightness(1.08);
}

.bulk-slot-actions :deep(.btn-primary:disabled),
.bulk-slot-actions :deep(.btn-neutral:disabled),
.bulk-slot-actions :deep(.btn-danger:disabled) {
  opacity: 0.5;
}
</style>
