import { defineStore } from 'pinia'
import { useSettingsStore } from '@/stores/settings'
import i18n from '@/locales'

export type NotificationType = 'update' | 'operation' | 'service' | 'error'
type FilterType = 'all' | NotificationType

export interface NotificationItem {
  id: string
  type: NotificationType
  icon: string
  title: string
  message: string
  at: number
  time: string
  unread: boolean
}

interface ToastItem {
  id: string
  icon: string
  message: string
}

const HISTORY_KEY = 'gobrew-notifications'
const FILTER_KEY = 'gobrew-notification-filter'
const SCROLL_KEY = 'gobrew-notification-scroll-top'
const MAX_HISTORY = 24

function nowLabel(ts: number) {
  const locale = i18n.global.locale.value === 'zh' ? 'zh-CN' : 'en-US'
  return new Intl.DateTimeFormat(locale, {
    hour: '2-digit',
    minute: '2-digit',
  }).format(ts)
}

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    open: false,
    filter: 'all' as FilterType,
    items: [] as NotificationItem[],
    toasts: [] as ToastItem[],
  }),
  getters: {
    unreadCount: (state) => state.items.filter(item => item.unread).length,
    filteredItems: (state) => state.filter === 'all'
      ? state.items
      : state.items.filter(item => item.type === state.filter),
  },
  actions: {
    load() {
      try {
        const parsed = JSON.parse(localStorage.getItem(HISTORY_KEY) || '[]')
        this.items = Array.isArray(parsed)
          ? parsed.map((item) => {
              const at = Number(item?.at || Date.now())
              return {
                ...item,
                at,
                time: item?.time || nowLabel(at),
              }
            })
          : []
      } catch {
        this.items = []
      }
      const savedFilter = localStorage.getItem(FILTER_KEY)
      if (savedFilter === 'all' || savedFilter === 'update' || savedFilter === 'operation' || savedFilter === 'service' || savedFilter === 'error') {
        this.filter = savedFilter
      }
    },
    save() {
      localStorage.setItem(HISTORY_KEY, JSON.stringify(this.items))
    },
    toggle() {
      this.open = !this.open
    },
    close() {
      this.open = false
    },
    setFilter(next: FilterType) {
      this.filter = next
      localStorage.setItem(FILTER_KEY, next)
    },
    saveScrollTop(top: number) {
      localStorage.setItem(SCROLL_KEY, String(Math.max(0, Math.floor(top))))
    },
    loadScrollTop() {
      const top = Number(localStorage.getItem(SCROLL_KEY) || '0')
      return Number.isFinite(top) ? Math.max(0, Math.floor(top)) : 0
    },
    allowed(type: NotificationType) {
      const settings = useSettingsStore()
      if (type === 'update') return settings.notifyUpdates
      if (type === 'service') return settings.notifyServices
      if (type === 'error') return settings.notifyErrors
      return settings.notifyOperations
    },
    push(type: NotificationType, icon: string, title: string, message: string, force = false) {
      const settings = useSettingsStore()
      if (!force && !this.allowed(type)) return
      if (settings.keepNotificationHistory) {
        const at = Date.now()
        this.items.unshift({
          id: `n-${at}-${Math.random().toString(36).slice(2, 6)}`,
          type,
          icon,
          title,
          message,
          at,
          time: nowLabel(at),
          unread: true,
        })
        this.items = this.items.slice(0, MAX_HISTORY)
        this.save()
      }
      if (!settings.showToasts) return
      this.toasts.push({
        id: `t-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
        icon,
        message: title,
      })
      const duration = Number(settings.toastDurationMs || 2500)
      const toastId = this.toasts[this.toasts.length - 1].id
      setTimeout(() => {
        this.toasts = this.toasts.filter((toast) => toast.id !== toastId)
      }, duration)
    },
    markRead(id: string) {
      this.items = this.items.map(item => item.id === id ? { ...item, unread: false } : item)
      this.save()
    },
    markAllRead() {
      this.items = this.items.map(item => ({ ...item, unread: false }))
      this.save()
    },
    clearAll() {
      this.items = []
      this.save()
    },
  },
})
