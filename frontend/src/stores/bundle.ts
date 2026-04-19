import { defineStore } from 'pinia'
import * as BundleService from '../../bindings/changeme/services/bundleservice.js'

interface BundleCheckResult {
  satisfied: boolean
  missing: string[]
  output: string
}

interface CommandResult {
  success: boolean
  output: string
  duration: string
  error?: string
}

export const useBundleStore = defineStore('bundle', {
  state: () => ({
    brewfileContent: '',
    checkResult: null as BundleCheckResult | null,
    cleanupPreviewResult: null as CommandResult | null,
    loading: false,
    error: '',
  }),
  actions: {
    async readBrewfile(path?: string) {
      this.loading = true
      this.error = ''
      try {
        this.brewfileContent = await BundleService.ReadBrewfile(path || '')
      } catch (error: any) {
        this.error = error?.message || 'Failed to read Brewfile'
        this.brewfileContent = ''
      } finally {
        this.loading = false
      }
    },
    async dump(path?: string, force = true) {
      await BundleService.Dump(path || '', force)
      await this.readBrewfile(path)
    },
    async restore(path?: string) {
      await BundleService.Restore(path || '')
    },
    async check(path?: string) {
      const result = await BundleService.Check(path || '')
      this.checkResult = (result as unknown) as BundleCheckResult
    },
    async cleanupPreview(path?: string) {
      const result = await BundleService.CleanupPreview(path || '')
      this.cleanupPreviewResult = (result as unknown) as CommandResult
    },
    async cleanup(path?: string, force = false) {
      await BundleService.Cleanup(path || '', force)
    },
  },
})
