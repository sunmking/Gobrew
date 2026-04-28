import { defineStore } from 'pinia'
import * as TapService from '../../bindings/changeme/services/tapservice.js'
import type { TapDetail, TapInfo } from '@/types/brew'
import i18n from '@/locales'

export const useTapsStore = defineStore('taps', {
  state: () => ({
    taps: [] as TapInfo[],
    selectedDetail: null as TapDetail | null,
    loading: false,
    error: '',
  }),
  actions: {
    async fetchTaps() {
      if (this.loading) return
      this.loading = true
      this.error = ''
      try {
        const result = await TapService.List()
        this.taps = (result as unknown) as TapInfo[]
      } catch (error: any) {
        this.error = error?.message || i18n.global.t('messages.failedFetchTaps')
      } finally {
        this.loading = false
      }
    },
    async add(name: string) {
      await TapService.Add(name)
      await this.fetchTaps()
    },
    async remove(name: string) {
      await TapService.Remove(name)
      await this.fetchTaps()
    },
    async details(name: string) {
      const result = await TapService.Details(name)
      this.selectedDetail = (result as unknown) as TapDetail
      return this.selectedDetail
    },
  },
})
