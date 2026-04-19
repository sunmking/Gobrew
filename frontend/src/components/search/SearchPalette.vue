<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSearch, type SearchResult } from '@/composables/useSearch'
import { Search, Package, Zap } from 'lucide-vue-next'

const { query, isOpen, results, close } = useSearch()
const { t } = useI18n()
const inputRef = ref<HTMLInputElement | null>(null)
const selectedIndex = ref(0)

const emit = defineEmits<{
  (e: 'select', result: SearchResult): void
}>()

watch(isOpen, async (open) => {
  if (open) {
    selectedIndex.value = 0
    await nextTick()
    inputRef.value?.focus()
  }
})

watch(results, () => {
  selectedIndex.value = 0
})

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, results.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (e.key === 'Enter' && results.value[selectedIndex.value]) {
    emit('select', results.value[selectedIndex.value])
    close()
  }
}

function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    if (isOpen.value) {
      close()
    } else {
      isOpen.value = true
      query.value = ''
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<template>
  <teleport to="body">
    <Transition name="palette">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh]" @click.self="close">
        <div class="fixed inset-0 bg-black/20 backdrop-blur-sm" @click="close" />
        <div class="relative w-full max-w-lg rounded-lg border shadow-lg bg-[var(--color-bg-elevated)] border-[var(--color-border)] overflow-hidden">
          <div class="flex items-center gap-2 px-3 py-2 border-b border-[var(--color-border)]">
            <Search :size="14" class="text-[var(--color-text-tertiary)] flex-shrink-0" />
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              :placeholder="t('searchPalette.placeholder')"
              class="flex-1 bg-transparent text-xs text-[var(--color-text-primary)] placeholder:text-[var(--color-text-tertiary)] outline-none"
              @keydown="handleKeydown"
            />
            <kbd class="text-[10px] text-[var(--color-text-tertiary)] bg-[var(--color-bg-muted)] px-1.5 py-0.5 rounded">ESC</kbd>
          </div>

          <div v-if="results.length > 0" class="max-h-64 overflow-y-auto py-1">
            <button
              v-for="(result, i) in results"
              :key="result.name + result.action"
              class="w-full flex items-center gap-2 px-3 py-1.5 text-left text-xs transition-colors"
              :class="i === selectedIndex ? 'bg-[var(--color-bg-muted)]' : 'hover:bg-[var(--color-bg-muted)]'"
              @click="emit('select', result); close()"
              @mouseenter="selectedIndex = i"
            >
              <component
                :is="result.type === 'package' ? Package : Zap"
                :size="14"
                class="flex-shrink-0 text-[var(--color-text-tertiary)]"
              />
              <div class="min-w-0 flex-1">
                <div class="text-[var(--color-text-primary)] font-medium truncate">{{ result.name }}</div>
                <div class="text-[var(--color-text-tertiary)] text-xs truncate">{{ result.description }}</div>
              </div>
            </button>
          </div>

          <div v-else-if="query.length > 0" class="py-8 text-center text-sm text-[var(--color-text-tertiary)]">
            {{ t('searchPalette.noResults') }}
          </div>

          <div v-else class="py-8 text-center text-sm text-[var(--color-text-tertiary)]">
            {{ t('searchPalette.typeHint') }}
          </div>

          <div class="flex items-center gap-3 px-4 py-2 border-t border-[var(--color-border)] text-[10px] text-[var(--color-text-tertiary)]">
            <span>{{ t('searchPalette.navigate') }}</span>
            <span>{{ t('searchPalette.select') }}</span>
            <span>{{ t('searchPalette.toggle') }}</span>
          </div>
        </div>
      </div>
    </Transition>
  </teleport>
</template>

<style scoped>
.palette-enter-active {
  transition: all 0.15s cubic-bezier(0.16, 1, 0.3, 1);
}
.palette-leave-active {
  transition: all 0.1s ease-in;
}
.palette-enter-from {
  opacity: 0;
  transform: scale(0.97) translateY(-8px);
}
.palette-leave-to {
  opacity: 0;
  transform: scale(0.97) translateY(-4px);
}
</style>
