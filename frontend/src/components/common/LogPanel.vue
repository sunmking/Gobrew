<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLogStore } from '@/stores/log'

const { t } = useI18n()

const logStore = useLogStore()
const collapsed = ref(true)
const scrollContainer = ref<HTMLElement | null>(null)

const recentLines = computed(() => logStore.lines.slice(-200))

const lineClass = (type: string) => {
  if (type === 'stderr') return 'text-red-400'
  if (type === 'system') return 'text-[#c9a962]'
  return 'text-green-400'
}

const scrollToBottom = async () => {
  await nextTick()
  if (scrollContainer.value) {
    scrollContainer.value.scrollTop = scrollContainer.value.scrollHeight
  }
}

watch(
  () => logStore.lines.length,
  () => {
    if (!collapsed.value) {
      scrollToBottom()
    }
  },
)

watch(
  () => logStore.listening,
  (running) => {
    if (running) {
      collapsed.value = false
    }
  },
)
</script>

<template>
  <div class="border-t border-gray-200 dark:border-[#222222] bg-gray-50 dark:bg-[#141414] text-green-300">
    <button
      class="flex w-full items-center gap-1.5 px-3 py-2 text-left text-xs text-gray-500 hover:text-gray-400 dark:text-[#666] dark:hover:text-[#999] hover:bg-gray-100 dark:hover:bg-[#1a1a1a]"
      @click="collapsed = !collapsed"
    >
      <span v-if="logStore.listening" class="inline-block h-2 w-2 animate-pulse rounded-full bg-[#c9a962]" />
      <span>{{ collapsed ? t('log.showLogs') : t('log.hideLogs') }} ({{ logStore.lines.length }})</span>
    </button>
    <div
      ref="scrollContainer"
      class="log-panel max-h-56 overflow-y-auto px-3 py-2 text-xs overflow-hidden transition-all duration-200"
      :class="collapsed ? 'max-h-0 !py-0 !px-3 opacity-0' : 'max-h-56 opacity-100'"
    >
      <div class="mb-2 flex justify-end">
        <button class="rounded bg-gray-200 text-gray-500 dark:bg-[#1a1a1a] dark:text-[#999] hover:bg-gray-300 dark:hover:bg-[#222] text-[11px] px-2 py-1" @click="logStore.clear()">
          {{ t('log.clear') }}
        </button>
      </div>
      <pre class="whitespace-pre-wrap break-words"><span v-for="line in recentLines" :key="line.id" :class="lineClass(line.type)">{{ line.text }}</span>&#10;</pre>
    </div>
  </div>
</template>
