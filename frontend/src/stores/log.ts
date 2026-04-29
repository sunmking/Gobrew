import { defineStore } from 'pinia'
import { Events } from '@wailsio/runtime'
import type { LogLine } from '@/types/brew'
import { useNotificationStore } from '@/stores/notification'
import i18n from '@/locales'

let lineId = 1
const ERROR_RE = /\b(error|failed|failure|panic)\b/i

export const useLogStore = defineStore('log', {
  state: () => ({
    lines: [] as LogLine[],
    listening: false,
    activeSessions: 0,
    stopFns: [] as Array<() => void>,
    maxLines: 1000,
    toastDurationMs: 2600,
    operationStatus: 'idle' as 'idle' | 'running' | 'success' | 'error' | 'info',
    operationMessage: '',
    clearTimer: null as ReturnType<typeof setTimeout> | null,
  }),
  actions: {
    setMaxLines(limit: number) {
      this.maxLines = Math.max(100, Math.floor(limit))
      if (this.lines.length > this.maxLines) {
        this.lines.splice(0, this.lines.length - this.maxLines)
      }
    },
    setToastDuration(ms: number) {
      this.toastDurationMs = Math.max(1000, Math.floor(ms))
    },
    setOperation(status: 'idle' | 'running' | 'success' | 'error' | 'info', message = '') {
      this.operationStatus = status
      this.operationMessage = message
    },
    scheduleClear(ms: number) {
      if (this.clearTimer) {
        clearTimeout(this.clearTimer)
      }
      this.clearTimer = setTimeout(() => {
        this.clearOperationStatus()
      }, ms)
    },
    markRunning(message = '') {
      if (this.clearTimer) {
        clearTimeout(this.clearTimer)
        this.clearTimer = null
      }
      this.setOperation('running', message)
    },
    markSuccess(message = '') {
      this.setOperation('success', message)
      this.scheduleClear(this.toastDurationMs)
      if (message) useNotificationStore().push('operation', '✓', message, i18n.global.t('notification.events.operationDoneMessage'))
    },
    markError(message = '') {
      this.setOperation('error', message)
      this.scheduleClear(Math.round(this.toastDurationMs * 2))
      if (message) useNotificationStore().push('error', '!', i18n.global.t('notification.events.operationErrorTitle'), message)
    },
    markInfo(message = '') {
      this.setOperation('info', message)
      this.scheduleClear(this.toastDurationMs)
      if (message) useNotificationStore().push('operation', 'ℹ', i18n.global.t('notification.events.operationInfoTitle'), message)
    },
    clearOperationStatus() {
      if (this.clearTimer) {
        clearTimeout(this.clearTimer)
        this.clearTimer = null
      }
      this.setOperation('idle', '')
    },
    startListening(message = '') {
      this.activeSessions += 1
      this.markRunning(message)

      if (this.listening) return
      this.listening = true

      const offOutput = Events.On('brew-output', (event) => {
        const text = String(event?.data ?? '')
        this.lines.push({
          id: lineId++,
          text,
          type: text.toLowerCase().includes('error') ? 'stderr' : 'stdout',
          timestamp: Date.now(),
        })
        if (this.lines.length > this.maxLines) {
          this.lines.splice(0, this.lines.length - this.maxLines)
        }
      })

      const offComplete = Events.On('brew-complete', (event) => {
        const text = `complete: ${String(event?.data ?? '')}`
        this.lines.push({
          id: lineId++,
          text,
          type: 'system',
          timestamp: Date.now(),
        })
        if (ERROR_RE.test(text)) {
          this.markError(text)
        } else {
          this.markSuccess(text)
        }
      })

      this.stopFns = [offOutput, offComplete]
    },
    stopListening() {
      this.activeSessions = Math.max(0, this.activeSessions - 1)
      if (this.activeSessions > 0) return

      this.stopFns.forEach((fn) => fn())
      this.stopFns = []
      this.listening = false
      if (this.operationStatus === 'running') {
        this.markSuccess('')
      }
    },
    clear() {
      this.lines = []
    },
  },
})
