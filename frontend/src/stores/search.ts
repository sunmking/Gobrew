import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import * as BrewService from '../../bindings/changeme/services/brewservice.js'
import type { SearchResult } from '@/types/brew'

export const useSearchStore = defineStore('search', () => {
  const results = reactive<SearchResult>({
    formulae: [],
    casks: [],
  })
  const query = ref('')
  const loading = ref(false)
  const error = ref('')

  function clearResults() {
    results.formulae = []
    results.casks = []
    error.value = ''
    query.value = ''
    loading.value = false
  }

  async function search(q: string) {
    results.formulae = []
    results.casks = []
    error.value = ''
    query.value = q
    if (!q.trim()) return

    loading.value = true
    try {
      const result = await BrewService.Search(q)
      const data = ((result ?? { formulae: [], casks: [] }) as unknown) as SearchResult
      results.formulae = data.formulae ?? []
      results.casks = data.casks ?? []
    } catch (err: any) {
      error.value = err?.message || 'Search failed'
    } finally {
      loading.value = false
    }
  }

  async function install(name: string) {
    await BrewService.Install(name)
  }

  return { results, query, loading, error, clearResults, search, install }
})
