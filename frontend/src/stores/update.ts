import { defineStore } from 'pinia'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import * as BundleService from '../../bindings/changeme/services/bundleservice.js'
import type { OutdatedCask, OutdatedFormula } from '@/types/brew'
import { useNotificationStore } from '@/stores/notification'
import { useSettingsStore } from '@/stores/settings'
import i18n from '@/locales'

export const useUpdateStore = defineStore('update', {
  state: () => ({
    formulae: [] as OutdatedFormula[],
    casks: [] as OutdatedCask[],
    loading: false,
    error: '',
    updating: false,
    loaded: false,
    lastNotifiedOutdated: -1,
  }),
  actions: {
    notifyOutdatedIfChanged() {
      const total = this.formulae.length + this.casks.length
      if (total === this.lastNotifiedOutdated) return
      this.lastNotifiedOutdated = total
      if (total <= 0) return
      const text = i18n.global.t('titleBar.updateFound', { count: total })
      useNotificationStore().push('update', '↻', text, text)
    },
    async fetchOutdated() {
      if (this.loaded && !this.loading) {
        this._silentRefresh()
        return
      }
      if (this.loading) return
      this.loading = true
      this.error = ''
      try {
        const result = await BrewService.Outdated()
        this.formulae = ((result?.formulae ?? []) as unknown) as OutdatedFormula[]
        this.casks = ((result?.casks ?? []) as unknown) as OutdatedCask[]
        this.loaded = true
        this.notifyOutdatedIfChanged()
      } catch (error: any) {
        this.error = error?.message || 'Failed to fetch outdated packages'
      } finally {
        this.loading = false
      }
    },
    async _silentRefresh() {
      try {
        const result = await BrewService.Outdated()
        this.formulae = ((result?.formulae ?? []) as unknown) as OutdatedFormula[]
        this.casks = ((result?.casks ?? []) as unknown) as OutdatedCask[]
        this.loaded = true
        this.notifyOutdatedIfChanged()
      } catch {
      }
    },
    async forceRefresh() {
      this.loading = true
      this.error = ''
      try {
        const result = await BrewService.Outdated()
        this.formulae = ((result?.formulae ?? []) as unknown) as OutdatedFormula[]
        this.casks = ((result?.casks ?? []) as unknown) as OutdatedCask[]
        this.loaded = true
        this.notifyOutdatedIfChanged()
      } catch (error: any) {
        this.error = error?.message || 'Failed to fetch outdated packages'
      } finally {
        this.loading = false
      }
    },
    async upgrade(name: string) {
      this.updating = true
      try {
        const settingsStore = useSettingsStore()
        if (settingsStore.backupBeforeUpdate) {
          await BundleService.Dump('', true)
        }
        await BrewService.Upgrade(name)
        if (settingsStore.cleanupAfterUpdate) {
          await BrewService.Cleanup()
        }
        await this.forceRefresh()
      } finally {
        this.updating = false
      }
    },
    async upgradeAll() {
      this.updating = true
      try {
        const settingsStore = useSettingsStore()
        if (settingsStore.backupBeforeUpdate) {
          await BundleService.Dump('', true)
        }
        await BrewService.UpgradeAll()
        if (settingsStore.cleanupAfterUpdate) {
          await BrewService.Cleanup()
        }
        await this.forceRefresh()
      } finally {
        this.updating = false
      }
    },
  },
})
