import { defineStore } from 'pinia'
import { Events } from '@wailsio/runtime'
import type { LogLine } from '@/types/brew'

let lineId = 1

export const useLogStore = defineStore('log', {
  state: () => ({
    lines: [] as LogLine[],
    listening: false,
    stopFns: [] as Array<() => void>,
    maxLines: 1000,
  }),
  actions: {
    startListening() {
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
      })

      this.stopFns = [offOutput, offComplete]
    },
    stopListening() {
      this.stopFns.forEach((fn) => fn())
      this.stopFns = []
      this.listening = false
    },
    clear() {
      this.lines = []
    },
  },
})
