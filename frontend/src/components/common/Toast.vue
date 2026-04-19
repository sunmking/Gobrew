<script setup lang="ts">
import { reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { CheckCircle, XCircle, Info, X, ChevronUp, ChevronDown } from 'lucide-vue-next'
import { useLogStore } from '@/stores/log'

type ToastType = 'success' | 'error' | 'info'

interface ToastItem {
  id: number
  type: ToastType
  message: string
  duration: number
  startTime: number
  expanded: boolean
  showProgress: boolean
}

const logStore = useLogStore()
const { t } = useI18n()
const toasts = reactive<ToastItem[]>([])
let nextId = 0
const MAX_TOASTS = 3

function dismiss(id: number) {
  const idx = toasts.findIndex(t => t.id === id)
  if (idx !== -1) toasts.splice(idx, 1)
}

function show(nextType: ToastType, nextMessage: string, duration = 3000, showProgress = false) {
  if (toasts.length >= MAX_TOASTS) toasts.shift()

  if (nextType === 'success') {
    logStore.markSuccess(nextMessage)
  } else if (nextType === 'error') {
    logStore.markError(nextMessage)
  } else {
    logStore.markInfo(nextMessage)
  }

  const id = nextId++
  toasts.push({
    id,
    type: nextType,
    message: nextMessage,
    duration,
    startTime: Date.now(),
    expanded: false,
    showProgress,
  })

  if (duration > 0) {
    setTimeout(() => dismiss(id), duration)
  }
}

function latestLogLine(): string {
  const lines = logStore.lines
  if (lines.length === 0) return ''
  return lines[lines.length - 1].text
}

function recentLogLines(): Array<{ text: string; type: string }> {
  return logStore.lines.slice(-50)
}

function toggleExpand(toast: ToastItem) {
  toast.expanded = !toast.expanded
}

const iconMap = {
  success: CheckCircle,
  error: XCircle,
  info: Info,
}

defineExpose({ show, dismiss })
</script>

<template>
  <teleport to="body">
    <div class="fixed bottom-4 right-4 z-50 flex flex-col-reverse gap-2 w-80 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="pointer-events-auto relative overflow-hidden rounded-md shadow-md text-sm"
          :class="{
            'toast-success': toast.type === 'success',
            'toast-error': toast.type === 'error',
            'toast-info': toast.type === 'info',
          }"
        >
          <div class="flex items-center gap-2 px-3 py-2.5">
            <component
              :is="iconMap[toast.type]"
              :size="16"
              :class="{
                'toast-icon-success': toast.type === 'success',
                'toast-icon-error': toast.type === 'error',
                'toast-icon-info': toast.type === 'info',
              }"
            />
            <span class="toast-message flex-1 text-xs">{{ toast.message }}</span>
            <button
              v-if="toast.showProgress && logStore.lines.length > 0"
              class="toast-expand-btn cursor-pointer"
              @click="toggleExpand(toast)"
            >
              <ChevronUp v-if="toast.expanded" :size="14" />
              <ChevronDown v-else :size="14" />
            </button>
            <button
              class="toast-dismiss-btn cursor-pointer"
              @click="dismiss(toast.id)"
            >
              <X :size="14" />
            </button>
          </div>

          <div
            v-if="toast.showProgress && !toast.expanded && latestLogLine()"
            class="toast-latest-line px-3 pb-2 text-[10px] font-mono truncate"
          >
            {{ latestLogLine() }}
          </div>

          <div
            v-if="toast.expanded"
            class="toast-log-panel max-h-40 overflow-y-auto px-3 py-2 font-mono text-[11px] whitespace-pre-wrap break-words"
          >
            <div v-for="(line, i) in recentLogLines()" :key="i" :class="{ 'toast-log-error': line.type === 'stderr' }">
              {{ line.text }}
            </div>
            <div v-if="recentLogLines().length === 0" class="toast-log-empty">{{ t('messages.noLogsYet') }}</div>
          </div>

          <div v-if="toast.duration > 0" class="toast-progress-track h-0.5 w-full">
            <div
              class="h-full"
              :class="{
                'toast-bar-success': toast.type === 'success',
                'toast-bar-error': toast.type === 'error',
                'toast-bar-info': toast.type === 'info',
              }"
              :style="{
                animation: `toast-progress ${toast.duration}ms linear forwards`,
              }"
            />
          </div>
        </div>
      </TransitionGroup>
    </div>
  </teleport>
</template>

<style scoped>
.toast-success {
  background: var(--color-card);
  border: 1px solid var(--color-success);
}

.toast-error {
  background: var(--color-card);
  border: 1px solid var(--color-danger);
}

.toast-info {
  background: var(--color-card);
  border: 1px solid var(--color-border);
}

.toast-icon-success {
  color: var(--color-success);
}

.toast-icon-error {
  color: var(--color-danger);
}

.toast-icon-info {
  color: var(--color-accent);
}

.toast-message {
  color: var(--color-text);
}

.toast-expand-btn,
.toast-dismiss-btn {
  color: var(--color-text-tertiary);
}

.toast-expand-btn:hover,
.toast-dismiss-btn:hover {
  color: var(--color-text-secondary);
}

.toast-latest-line {
  color: var(--color-text-tertiary);
}

.toast-log-panel {
  border-top: 1px solid var(--color-border);
  background: var(--color-group-bg);
  color: var(--color-text-secondary);
}

.toast-log-error {
  color: var(--color-danger);
}

.toast-log-empty {
  color: var(--color-text-tertiary);
}

.toast-progress-track {
  background: var(--color-group-bg);
}

.toast-bar-success {
  background: var(--color-success);
}

.toast-bar-error {
  background: var(--color-danger);
}

.toast-bar-info {
  background: var(--color-accent);
}

.toast-enter-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-leave-active {
  transition: all 0.2s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.95);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.95);
}
.toast-move {
  transition: transform 0.2s ease;
}

@keyframes toast-progress {
  from { width: 100%; }
  to { width: 0%; }
}
</style>
