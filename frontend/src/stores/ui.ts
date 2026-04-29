import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    settingsOpen: false,
  }),
  actions: {
    openSettings() {
      this.settingsOpen = true
    },
    closeSettings() {
      this.settingsOpen = false
    },
    toggleSettings() {
      this.settingsOpen = !this.settingsOpen
    },
  },
})
