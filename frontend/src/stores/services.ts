import { defineStore } from 'pinia'
import * as ServiceManager from '../../bindings/changeme/services/servicemanager.js'
import type { BrewService } from '@/types/brew'
import i18n from '@/locales'

export const useServicesStore = defineStore('services', {
  state: () => ({
    services: [] as BrewService[],
    loading: false,
    error: '',
  }),
  actions: {
    async fetchServices() {
      if (this.loading) return
      this.loading = true
      this.error = ''
      try {
        const result = await ServiceManager.List()
        this.services = (result as unknown) as BrewService[]
      } catch (error: any) {
        this.error = error?.message || i18n.global.t('messages.failedFetchServices')
      } finally {
        this.loading = false
      }
    },
    async start(name: string) {
      await ServiceManager.Start(name)
      await this.fetchServices()
    },
    async stop(name: string) {
      await ServiceManager.Stop(name)
      await this.fetchServices()
    },
    async restart(name: string) {
      await ServiceManager.Restart(name)
      await this.fetchServices()
    },
    async startAll() {
      await ServiceManager.StartAll()
      await this.fetchServices()
    },
    async stopAll() {
      await ServiceManager.StopAll()
      await this.fetchServices()
    },
    async restartAll() {
      await ServiceManager.RestartAll()
      await this.fetchServices()
    },
    async cleanup() {
      await ServiceManager.Cleanup()
      await this.fetchServices()
    },
  },
})
