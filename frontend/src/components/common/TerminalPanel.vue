<script setup lang="ts">
import type { LogLine } from '@/types/brew'

defineProps<{
  lines: Array<LogLine | string>
  empty?: string
}>()

function lineText(line: LogLine | string) {
  return typeof line === 'string' ? line : line.text
}

function lineClass(line: LogLine | string) {
  if (typeof line === 'string') return 't-out'
  if (line.type === 'system') return 't-ok'
  if (line.type === 'stderr') return 't-err'
  return 't-out'
}
</script>

<template>
  <div class="terminal">
    <div v-if="lines.length === 0" class="t-out">{{ empty || '暂无命令输出' }}</div>
    <div v-for="(line, index) in lines" :key="typeof line === 'string' ? `${index}-${line}` : line.id" :class="lineClass(line)">
      {{ lineText(line) }}
    </div>
  </div>
</template>
