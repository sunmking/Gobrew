import { defineStore } from 'pinia'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import type { CaskInstalled, FormulaInstalled } from '@/types/brew'
import i18n from '@/locales'

export const useInstalledStore = defineStore('installed', {
  state: () => ({
    formulae: [] as FormulaInstalled[],
    casks: [] as CaskInstalled[],
    loading: false,
    error: '',
    loaded: false,
  }),
  actions: {
    async fetchInstalled() {
      if (this.loaded && !this.loading) {
        this._silentRefresh()
        return
      }
      if (this.loading) return
      this.loading = true
      this.error = ''
      try {
        const result = await BrewService.ListInstalled()
        this.formulae = ((result?.formulae ?? []) as unknown) as FormulaInstalled[]
        this.casks = ((result?.casks ?? []) as unknown) as CaskInstalled[]
        this.loaded = true
      } catch (error: any) {
        this.error = error?.message || i18n.global.t('messages.failedFetchInstalled')
      } finally {
        this.loading = false
      }
    },
    async _silentRefresh() {
      try {
        const result = await BrewService.ListInstalled()
        this.formulae = ((result?.formulae ?? []) as unknown) as FormulaInstalled[]
        this.casks = ((result?.casks ?? []) as unknown) as CaskInstalled[]
        this.loaded = true
      } catch {
      }
    },
    async forceRefresh() {
      this.loading = true
      this.error = ''
      try {
        const result = await BrewService.ListInstalled()
        this.formulae = ((result?.formulae ?? []) as unknown) as FormulaInstalled[]
        this.casks = ((result?.casks ?? []) as unknown) as CaskInstalled[]
        this.loaded = true
      } catch (error: any) {
        this.error = error?.message || i18n.global.t('messages.failedFetchInstalled')
      } finally {
        this.loading = false
      }
    },
    async uninstall(name: string) {
      await BrewService.Uninstall(name)
      await this.forceRefresh()
    },
    async uninstallCask(name: string, zap = false) {
      await BrewService.UninstallCask(name, zap)
      await this.forceRefresh()
    },
    async upgrade(name: string) {
      await BrewService.Upgrade(name)
      await this.forceRefresh()
    },
    async reinstall(name: string, isCask = false) {
      await BrewService.Reinstall(name, isCask)
      await this.forceRefresh()
    },
  },
})
