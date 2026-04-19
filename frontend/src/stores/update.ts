import { defineStore } from 'pinia'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import type { OutdatedCask, OutdatedFormula } from '@/types/brew'

export const useUpdateStore = defineStore('update', {
  state: () => ({
    formulae: [] as OutdatedFormula[],
    casks: [] as OutdatedCask[],
    loading: false,
    error: '',
    updating: false,
    loaded: false,
  }),
  actions: {
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
      } catch (error: any) {
        this.error = error?.message || 'Failed to fetch outdated packages'
      } finally {
        this.loading = false
      }
    },
    async upgrade(name: string) {
      this.updating = true
      try {
        await BrewService.Upgrade(name)
        await this.forceRefresh()
      } finally {
        this.updating = false
      }
    },
    async upgradeAll() {
      this.updating = true
      try {
        await BrewService.UpgradeAll()
        await this.forceRefresh()
      } finally {
        this.updating = false
      }
    },
  },
})
