<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import SettingsPage from '@/pages/SettingsPage.vue'
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
const { t } = useI18n()

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    uiStore.closeSettings()
  }
}

function onOverlayClick(event: MouseEvent) {
  if ((event.target as HTMLElement)?.classList.contains('prefs-overlay')) {
    uiStore.closeSettings()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <teleport to="body">
    <div v-if="uiStore.settingsOpen" class="prefs-overlay" @click="onOverlayClick">
      <section class="prefs-window" role="dialog" aria-modal="true" :aria-label="t('settingsPrefs.title')">
        <button class="prefs-close" type="button" @click="uiStore.closeSettings">×</button>
        <SettingsPage />
      </section>
    </div>
  </teleport>
</template>

<style scoped>
.prefs-overlay {
  position: fixed;
  inset: 0;
  z-index: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  background: oklch(0% 0 0 / 0.42);
  backdrop-filter: blur(8px);
}
.prefs-window {
  width: min(660px, calc(100vw - 24px));
  height: min(508px, calc(100vh - 24px));
  background: var(--bg);
  border: 1px solid oklch(0% 0 0 / 0.22);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  position: relative;
  box-shadow: 0 28px 80px oklch(0% 0 0 / 0.48);
}
.prefs-close {
  position: absolute;
  top: 10px;
  right: 12px;
  z-index: 10;
  width: 26px;
  height: 26px;
  border: 0;
  border-radius: 50%;
  background: var(--border);
  color: var(--muted);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
}
.prefs-close:hover { background: var(--muted); color: var(--bg); }
.prefs-close:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

@media (max-width: 820px) {
  .prefs-window {
    border-radius: 10px;
  }
}
</style>
