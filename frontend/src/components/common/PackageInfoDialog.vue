<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingInline from '@/components/common/LoadingInline.vue'

export interface InfoRow {
  label: string
  value: string
}

defineProps<{
  open: boolean
  title: string
  rows: InfoRow[]
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()
</script>

<template>
  <teleport to="body">
    <div v-if="open" class="info-backdrop fixed inset-0 z-40 flex items-center justify-center p-4">
      <div class="info-panel w-full max-w-xl rounded-xl p-4 shadow-xl">
        <div class="flex items-center justify-between gap-2">
          <h3 class="info-title truncate text-base font-semibold">{{ title }}</h3>
          <button class="info-close rounded px-2 py-1 text-xs" @click="emit('close')">{{ t('common.close') }}</button>
        </div>

        <div class="info-content mt-3 max-h-[60vh] overflow-y-auto rounded-lg p-3">
          <div v-if="loading" class="info-muted text-sm"><LoadingInline compact /></div>
          <div v-else-if="error" class="info-error rounded-md p-2 text-sm">{{ error }}</div>
          <div v-else-if="rows.length === 0" class="info-muted text-sm">{{ t('common.noDetails') }}</div>
          <dl v-else class="space-y-2">
            <div v-for="row in rows" :key="row.label" class="grid grid-cols-3 gap-2 text-sm">
              <dt class="info-label">{{ row.label }}</dt>
              <dd class="info-value col-span-2 whitespace-pre-wrap break-words">{{ row.value || '-' }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  </teleport>
</template>

<style scoped>
.info-backdrop {
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.info-panel {
  background: var(--color-group-bg);
  border: 1px solid var(--color-border);
}

.info-title {
  color: var(--color-text);
}

.info-close {
  background: var(--color-group-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.info-close:hover {
  opacity: 0.85;
}

.info-content {
  background: var(--color-group-bg);
  border: 1px solid var(--color-border);
}

.info-muted {
  color: var(--color-text-tertiary);
}

.info-error {
  background: var(--color-danger-light);
  border: 1px solid var(--color-danger);
  color: var(--color-danger);
}

.info-label {
  color: var(--color-text-secondary);
}

.info-value {
  color: var(--color-text);
}
</style>
