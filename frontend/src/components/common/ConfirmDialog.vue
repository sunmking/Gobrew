<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  TransitionRoot,
  TransitionChild,
} from '@headlessui/vue'

const props = defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  severity?: 'default' | 'danger'
  danger?: boolean
  toggleLabel?: string
  toggleValue?: boolean
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
  'update:toggleValue': [value: boolean]
}>()

const { t } = useI18n()
const effectiveSeverity = computed(() => props.severity ?? (props.danger ? 'danger' : 'default'))

function onToggleChange(event: Event) {
  emit('update:toggleValue', (event.target as HTMLInputElement).checked)
}
</script>

<template>
  <TransitionRoot appear :show="open" as="template">
    <Dialog as="div" class="relative z-40" @close="emit('cancel')">
      <TransitionChild
        as="template"
        enter="duration-200 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-150 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="confirm-backdrop fixed inset-0" />
      </TransitionChild>

      <div class="fixed inset-0 flex items-center justify-center p-4">
        <TransitionChild
          as="template"
          enter="duration-200 ease-out"
          enter-from="opacity-0 scale-95"
          enter-to="opacity-100 scale-100"
          leave="duration-150 ease-in"
          leave-from="opacity-100 scale-100"
          leave-to="opacity-0 scale-95"
        >
          <DialogPanel class="confirm-panel w-full max-w-md rounded-xl p-4 shadow-xl">
            <DialogTitle class="confirm-title text-base font-semibold">
              {{ title }}
            </DialogTitle>
            <p class="confirm-message mt-2 text-sm">{{ message }}</p>
            <label v-if="toggleLabel" class="confirm-toggle mt-3 flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                class="h-4 w-4"
                :checked="!!toggleValue"
                @change="onToggleChange"
              />
              {{ toggleLabel }}
            </label>
            <div class="mt-4 flex justify-end gap-2">
              <button
                class="confirm-cancel rounded-md px-3 py-1.5 text-sm"
                @click="emit('cancel')"
              >
                {{ cancelLabel || t('common.cancel') }}
              </button>
              <button
                class="confirm-action rounded-md px-4 py-2 text-sm"
                :class="{ 'confirm-action--danger': effectiveSeverity === 'danger' }"
                @click="emit('confirm')"
              >
                {{ confirmLabel || t('common.confirm') }}
              </button>
            </div>
          </DialogPanel>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<style scoped>
.confirm-backdrop {
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.confirm-panel {
  background: var(--color-group-bg);
  border: 1px solid var(--color-border);
}

.confirm-title {
  color: var(--color-text);
}

.confirm-message,
.confirm-toggle {
  color: var(--color-text-secondary);
}

.confirm-cancel {
  background: var(--color-group-bg);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.confirm-cancel:hover {
  opacity: 0.85;
}

.confirm-action {
  background: var(--color-accent);
  color: white;
  font-weight: 500;
}

.confirm-action:hover {
  background: var(--color-accent-hover);
}

.confirm-action--danger {
  background: var(--color-danger);
}

.confirm-action--danger:hover {
  filter: brightness(1.08);
}
</style>
