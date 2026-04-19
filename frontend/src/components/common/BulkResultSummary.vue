<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { BulkSummary } from '@/types/bulk'

defineProps<{
  summary: BulkSummary | null
}>()

const expanded = ref(false)
const { t } = useI18n()
</script>

<template>
  <div v-if="summary" class="summary-container mb-3 rounded-lg p-3 text-xs">
    <div class="font-semibold" style="color: var(--color-text)">
      {{ summary.action }}
    </div>
    <div class="mt-2 grid grid-cols-3 gap-2 text-[11px]">
      <div class="summary-cell summary-cell--total rounded px-2 py-1">
        {{ t('bulk.total') }}: <span class="font-semibold">{{ summary.total }}</span>
      </div>
      <div class="summary-cell summary-cell--success rounded px-2 py-1">
        {{ t('bulk.succeeded') }}: <span class="font-semibold">{{ summary.success }}</span>
      </div>
      <div class="summary-cell summary-cell--failure rounded px-2 py-1">
        {{ t('bulk.failed') }}: <span class="font-semibold">{{ summary.failures.length }}</span>
      </div>
    </div>
    <div v-if="summary.failures.length > 0" class="mt-3">
      <button class="summary-expand font-medium hover:underline" @click="expanded = !expanded">
        {{ expanded ? t('bulk.hideDetails') : t('bulk.showDetails') }}
      </button>
      <ul
        v-if="expanded"
        class="summary-failure-list mt-2 space-y-1 rounded p-2"
      >
        <li v-for="failure in summary.failures" :key="`${failure.item}-${failure.reason}`">
          {{ failure.item }} - {{ failure.reason || t('bulk.reasonUnknown') }}
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.summary-container {
  background: var(--color-card);
  border: 1px solid var(--color-border);
}

.summary-cell {
  border: 1px solid var(--color-border);
}

.summary-cell--total {
  color: var(--color-text-secondary);
}

.summary-cell--success {
  background: var(--color-success-light);
  color: var(--color-success);
  border-color: var(--color-success);
}

.summary-cell--failure {
  background: var(--color-danger-light);
  color: var(--color-danger);
  border-color: var(--color-danger);
}

.summary-expand {
  color: var(--color-danger);
}

.summary-failure-list {
  background: var(--color-danger-light);
  border: 1px solid var(--color-danger);
  color: var(--color-danger);
}
</style>
