<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  name: string
  size?: number
}>()

const size = computed(() => props.size ?? 28)

const builtinIcons: Record<string, string> = {
  node: '#339933',
  python: '#3776AB',
  pip: '#3776AB',
  go: '#00ADD8',
  rust: '#CE422B',
  java: '#ED8B00',
  javascript: '#F7DF1E',
  typescript: '#3178C6',
  docker: '#2496ED',
  git: '#F05032',
  vim: '#019833',
  nvim: '#57A143',
  nginx: '#009639',
  postgres: '#336791',
  mysql: '#4479A1',
  redis: '#DC382D',
  yarn: '#2C8EBB',
  pnpm: '#F69220',
  composer: '#885630',
  ruby: '#CC342D',
  php: '#777BB4',
  swift: '#FA7343',
  kotlin: '#7F52FF',
  cmake: '#064F8C',
  gradle: '#02303A',
  deno: '#000000',
  brew: '#0a0a0a',
  homebrew: '#0a0a0a',
  wget: '#6a5acd',
  curl: '#073551',
  tree: '#4CAF50',
  jq: '#0e8b93',
  htop: '#4CAF50',
  terraform: '#7B42BC',
  kubectl: '#326CE5',
  helm: '#0F1689',
  cocoapods: '#EE3322',
  npm: '#CB3837',
  cargo: '#DEA584',
  bundler: '#CC342D',
}

const isBuiltin = computed(() => props.name.toLowerCase() in builtinIcons)

const bgColor = computed(() => {
  if (isBuiltin.value) return builtinIcons[props.name.toLowerCase()] || '#636366'
  const name = props.name.toLowerCase().replace(/[^a-z0-9]/g, '')
  const colors = [
    '#4A90D9', '#D35D6E', '#50B83C', '#E8893C', '#7B68EE',
    '#2AA198', '#CF6A87', '#5C6BC0', '#26A69A', '#8D6E63',
    '#7E57C2', '#42A5F5', '#EF5350', '#66BB6A', '#FFA726',
    '#AB47BC', '#5C6BC0', '#EC407A', '#26C6DA', '#FFCA28',
  ]
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
})

const letter = computed(() => {
  const name = props.name.replace(/[^a-zA-Z0-9]/g, '')
  return name ? name.charAt(0).toUpperCase() : '?'
})
</script>

<template>
  <div
    :style="{
      width: size + 'px',
      height: size + 'px',
      minWidth: size + 'px',
      minHeight: size + 'px',
      borderRadius: '6px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: bgColor,
      flexShrink: 0,
      fontSize: (size * 0.42) + 'px',
      fontWeight: 600,
      fontFamily: '-apple-system, BlinkMacSystemFont, sans-serif',
      color: 'white',
      letterSpacing: '-0.3px',
      lineHeight: 1,
      overflow: 'hidden',
    }"
  >
    <span v-if="!isBuiltin">{{ letter }}</span>
  </div>
</template>
