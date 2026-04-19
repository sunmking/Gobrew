<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Loader2, CheckCircle2, AlertCircle, Info, X, ChevronUp, ChevronDown } from 'lucide-vue-next'
import { useLogStore } from '@/stores/log'

const { t } = useI18n()
const logStore = useLogStore()
const expanded = ref(false)

const visible = computed(() => logStore.operationStatus !== 'idle')
const statusClass = computed(() => `is-${logStore.operationStatus}`)
const statusText = computed(() => {
  if (logStore.operationStatus === 'running') return t('common.operationRunning')
  if (logStore.operationStatus === 'success') return t('common.operationSuccess')
  if (logStore.operationStatus === 'error') return t('common.operationError')
  return t('common.operationInfo')
})
const detailText = computed(() => {
  if (logStore.operationMessage) return logStore.operationMessage
  const line = logStore.lines[logStore.lines.length - 1]
  return line?.text || t('common.operationNoDetails')
})
const showLogs = computed(() => logStore.lines.length > 0)

function closeBar() {
  expanded.value = false
  logStore.clearOperationStatus()
}
</script>

<template>
  <Transition name="opbar">
    <div v-if="visible" class="operation-bar" :class="statusClass">
      <div class="operation-main">
        <div class="operation-leading">
          <Loader2 v-if="logStore.operationStatus === 'running'" :size="16" class="spinning" />
          <CheckCircle2 v-else-if="logStore.operationStatus === 'success'" :size="16" />
          <AlertCircle v-else-if="logStore.operationStatus === 'error'" :size="16" />
          <Info v-else :size="16" />
        </div>
        <div class="operation-content">
          <div class="operation-title">{{ statusText }}</div>
          <div class="operation-detail">{{ detailText }}</div>
        </div>
        <button
          v-if="showLogs"
          class="operation-btn"
          @click="expanded = !expanded"
        >
          <ChevronUp v-if="expanded" :size="14" />
          <ChevronDown v-else :size="14" />
          <span>{{ expanded ? t('log.hideLogs') : t('log.showLogs') }}</span>
        </button>
        <button class="operation-btn" @click="closeBar">
          <X :size="14" />
        </button>
      </div>
      <div v-if="expanded" class="operation-logs">
        <div
          v-for="line in logStore.lines.slice(-30)"
          :key="line.id"
          class="operation-log-line"
          :class="{ 'is-error': line.type === 'stderr' }"
        >
          {{ line.text }}
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.operation-bar {
  position: fixed;
  left: 236px;
  right: 16px;
  bottom: 14px;
  z-index: 46;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--color-card) 96%, transparent);
  backdrop-filter: blur(10px);
  box-shadow: var(--color-shadow);
}

.operation-main {
  min-height: 40px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
}

.operation-leading {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.operation-content {
  flex: 1;
  min-width: 0;
}

.operation-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
}

.operation-detail {
  margin-top: 2px;
  font-size: 11px;
  color: var(--color-text-secondary);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.operation-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  padding: 4px 6px;
  font-size: 11px;
  cursor: pointer;
}

.operation-btn:hover {
  background: var(--color-sidebar-hover);
}

.operation-logs {
  border-top: 1px solid var(--color-border);
  max-height: 160px;
  overflow: auto;
  padding: 8px 10px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--color-text-secondary);
  background: var(--color-group-bg);
}

.operation-log-line + .operation-log-line {
  margin-top: 4px;
}

.operation-log-line.is-error {
  color: var(--color-danger);
}

.is-running .operation-leading {
  color: var(--color-accent);
}

.is-success .operation-leading {
  color: var(--color-success);
}

.is-error .operation-leading {
  color: var(--color-danger);
}

.is-info .operation-leading {
  color: var(--color-text-secondary);
}

.spinning {
  animation: spin 1s linear infinite;
}

.opbar-enter-active,
.opbar-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.opbar-enter-from,
.opbar-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (max-width: 920px) {
  .operation-bar {
    left: 12px;
    right: 12px;
  }
}
</style>
