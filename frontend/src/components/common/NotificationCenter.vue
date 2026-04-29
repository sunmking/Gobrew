<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '@/stores/notification'
import { useSettingsStore } from '@/stores/settings'

const { t } = useI18n()
const notificationStore = useNotificationStore()
const settingsStore = useSettingsStore()
const panelRef = ref<HTMLElement | null>(null)
const listRef = ref<HTMLElement | null>(null)

const panelClass = computed(() => ({ open: notificationStore.open }))
const toasterStyle = computed(() => {
  const pos = settingsStore.toastPosition
  return {
    top: pos.startsWith('t') ? '20px' : 'auto',
    bottom: pos.startsWith('b') ? '20px' : 'auto',
    left: pos.endsWith('l') ? '20px' : 'auto',
    right: pos.endsWith('r') ? '20px' : 'auto',
  }
})

function iconClass(type: string) {
  return `notif-icon ${type}`
}

function onDocumentClick(event: MouseEvent) {
  const target = event.target as HTMLElement | null
  if (!target) return
  if (target.closest('.notif-button')) return
  if (panelRef.value?.contains(target)) return
  notificationStore.close()
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    notificationStore.close()
  }
}

onMounted(() => {
  notificationStore.load()
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
})

watch(
  () => notificationStore.open,
  async (open) => {
    if (!open) return
    await nextTick()
    const top = notificationStore.loadScrollTop()
    if (listRef.value) listRef.value.scrollTop = top
    const firstFocusable = panelRef.value?.querySelector<HTMLElement>('.notif-filter, .notif-item')
    firstFocusable?.focus()
  },
)

function onListScroll() {
  if (!listRef.value) return
  notificationStore.saveScrollTop(listRef.value.scrollTop)
}
</script>

<template>
  <div ref="panelRef" class="notification-panel" :class="panelClass" @click.stop>
    <div class="notif-head">
      <div class="notif-title">{{ t('notification.title') }}</div>
      <div class="notif-actions">
        <button class="btn btn-ghost" style="height:26px;font-size:12px;" @click="notificationStore.clearAll()">{{ t('notification.clear') }}</button>
        <button class="btn btn-ghost" style="height:26px;font-size:12px;" @click="notificationStore.markAllRead()">{{ t('notification.markAllRead') }}</button>
        <button class="btn-icon" :title="t('common.close')" @click="notificationStore.close()">×</button>
      </div>
    </div>
    <div class="notif-filters">
      <button class="notif-filter" :class="{ active: notificationStore.filter === 'all' }" @click="notificationStore.setFilter('all')">{{ t('notification.filters.all') }}</button>
      <button class="notif-filter" :class="{ active: notificationStore.filter === 'update' }" @click="notificationStore.setFilter('update')">{{ t('notification.filters.update') }}</button>
      <button class="notif-filter" :class="{ active: notificationStore.filter === 'operation' }" @click="notificationStore.setFilter('operation')">{{ t('notification.filters.operation') }}</button>
      <button class="notif-filter" :class="{ active: notificationStore.filter === 'service' }" @click="notificationStore.setFilter('service')">{{ t('notification.filters.service') }}</button>
      <button class="notif-filter" :class="{ active: notificationStore.filter === 'error' }" @click="notificationStore.setFilter('error')">{{ t('notification.filters.error') }}</button>
    </div>
    <div ref="listRef" class="notif-list" @scroll="onListScroll">
      <div v-if="notificationStore.filteredItems.length === 0" class="notif-empty">{{ t('notification.empty') }}</div>
      <button
        v-for="item in notificationStore.filteredItems"
        :key="item.id"
        class="notif-item"
        type="button"
        :class="{ unread: item.unread }"
        @click="notificationStore.markRead(item.id)"
      >
        <div :class="iconClass(item.type)">{{ item.icon }}</div>
        <div class="notif-content">
          <div class="notif-name">{{ item.title }}</div>
          <div class="notif-message">{{ item.message }}</div>
        </div>
        <div class="notif-time">{{ item.time }}</div>
      </button>
    </div>
  </div>

  <teleport to="body">
    <div class="toaster" :style="toasterStyle">
      <div v-for="toast in notificationStore.toasts" :key="toast.id" class="toast">
        <span>{{ toast.icon }}</span><span>{{ toast.message }}</span>
      </div>
    </div>
  </teleport>
</template>

<style scoped>
.toaster { position:fixed; display:flex; flex-direction:column; gap:8px; z-index:600; pointer-events:none; }
.toast { background:var(--fg); color:var(--bg); padding:10px 14px; border-radius:var(--radius); font-size:13px; font-weight:450; animation:slide-in 0.2s ease; max-width:300px; display:flex; align-items:center; gap:8px; }
@keyframes slide-in { from { transform:translateX(20px); opacity:0; } to { transform:translateX(0); opacity:1; } }
.notification-panel {
  position:absolute; top:42px; right:48px; z-index:450;
  width:min(360px, calc(100vw - 24px)); max-height:min(520px, calc(100vh - 72px));
  background:var(--surface); color:var(--fg);
  border:1px solid var(--border); border-radius:12px;
  box-shadow:0 18px 60px oklch(0% 0 0 / 0.26), 0 0 0 0.5px oklch(0% 0 0 / 0.12);
  overflow:hidden;
  opacity:0; transform:translateY(-6px) scale(0.98);
  pointer-events:none;
  transition:opacity 0.16s ease, transform 0.16s ease;
}
.notification-panel.open { opacity:1; transform:translateY(0) scale(1); pointer-events:all; }
.notification-panel::before {
  content:""; position:absolute; top:-6px; right:22px;
  width:12px; height:12px; background:var(--surface);
  border-left:1px solid var(--border); border-top:1px solid var(--border);
  transform:rotate(45deg);
}
.notif-head { display:flex; align-items:center; justify-content:space-between; padding:14px 14px 10px; border-bottom:1px solid var(--border); }
.notif-title { font-size:14px; font-weight:650; letter-spacing:-0.015em; }
.notif-actions { display:flex; gap:4px; }
.notif-filters { display:flex; gap:4px; padding:8px 10px; border-bottom:1px solid var(--border); overflow-x:auto; }
.notif-filter { border:0; background:transparent; color:var(--muted); height:26px; padding:0 8px; border-radius:var(--radius-sm); font:500 12px var(--font-body); cursor:pointer; white-space:nowrap; }
.notif-filter:hover { background:var(--border); color:var(--fg); }
.notif-filter.active { background:var(--accent-bg); color:var(--accent); }
.notif-list { max-height:360px; overflow:auto; }
.notif-list { scrollbar-gutter: stable; }
.notif-item { width:100%; border:0; text-align:left; display:grid; grid-template-columns:26px 1fr auto; gap:10px; padding:12px 14px; border-bottom:1px solid var(--border); background:transparent; }
.notif-item { cursor: pointer; transition: background 0.12s ease; }
.notif-item:last-child { border-bottom:0; }
.notif-item.unread { background:color-mix(in oklch, var(--accent-bg) 38%, transparent); }
.notif-item:hover { background: color-mix(in oklch, var(--border) 50%, transparent); }
.notif-icon { width:26px; height:26px; border-radius:50%; display:flex; align-items:center; justify-content:center; background:var(--border); color:var(--fg); font-size:12px; font-weight:650; }
.notif-icon.update { background:oklch(93% 0.04 255); color:oklch(38% 0.18 255); }
.notif-icon.operation { background:oklch(93% 0.04 145); color:oklch(38% 0.14 145); }
.notif-icon.service { background:oklch(94% 0.04 70); color:oklch(42% 0.13 70); }
.notif-icon.error { background:oklch(94% 0.04 25); color:oklch(40% 0.18 25); }
[data-theme="dark"] .notif-icon.update { background:oklch(25% 0.08 255); color:oklch(72% 0.18 255); }
[data-theme="dark"] .notif-icon.operation { background:oklch(25% 0.08 145); color:oklch(72% 0.14 145); }
[data-theme="dark"] .notif-icon.service { background:oklch(25% 0.08 70); color:oklch(78% 0.13 70); }
[data-theme="dark"] .notif-icon.error { background:oklch(25% 0.08 25); color:oklch(72% 0.18 25); }
.notif-content { min-width:0; }
.notif-name { font-size:12.8px; font-weight:550; color:var(--fg); }
.notif-message { margin-top:2px; font-size:11.8px; line-height:1.45; color:var(--muted); }
.notif-time { font-family:var(--font-mono); font-size:10.5px; color:var(--muted); white-space:nowrap; }
.notif-empty { padding:36px 18px; text-align:center; color:var(--muted); font-size:12.5px; }
.notif-filter:focus-visible,
.notif-item:focus-visible,
.notif-actions .btn:focus-visible,
.notif-actions .btn-icon:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 1px;
}

@media (max-width: 820px) {
  .notification-panel {
    right: 12px;
    top: 40px;
  }
}

@media (prefers-reduced-motion: no-preference) {
  .notif-item {
    animation: notif-fade-in 0.16s ease;
  }
  @keyframes notif-fade-in {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
}
</style>
