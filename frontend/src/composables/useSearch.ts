import { ref, computed } from 'vue'
import { useInstalledStore } from '@/stores/installed'
import { useUpdateStore } from '@/stores/update'

export interface SearchResult {
  type: 'package' | 'action'
  name: string
  description: string
  action: string
}

export function useSearch() {
  const query = ref('')
  const isOpen = ref(false)

  const installedStore = useInstalledStore()
  const updateStore = useUpdateStore()

  function open() {
    isOpen.value = true
    query.value = ''
  }

  function close() {
    isOpen.value = false
    query.value = ''
  }

  const results = computed<SearchResult[]>(() => {
    const q = query.value.toLowerCase().trim()
    if (!q) return []

    const matches: SearchResult[] = []

    for (const pkg of installedStore.formulae) {
      if (pkg.name.toLowerCase().includes(q)) {
        matches.push({
          type: 'package',
          name: pkg.name,
          description: `Installed · ${pkg.installed?.[0]?.version || 'unknown'}`,
          action: 'info',
        })
      }
    }

    for (const pkg of installedStore.casks) {
      if (pkg.name.toLowerCase().includes(q)) {
        matches.push({
          type: 'package',
          name: pkg.name,
          description: `Installed (cask) · ${pkg.installed || 'unknown'}`,
          action: 'info',
        })
      }
    }

    for (const pkg of updateStore.formulae) {
      if (pkg.name.toLowerCase().includes(q)) {
        matches.push({
          type: 'package',
          name: pkg.name,
          description: `Upgradable · ${pkg.installed_versions?.[0] || 'unknown'} → ${pkg.current_version}`,
          action: 'upgrade',
        })
      }
    }

    for (const pkg of updateStore.casks) {
      if (pkg.name.toLowerCase().includes(q)) {
        matches.push({
          type: 'package',
          name: pkg.name,
          description: `Upgradable (cask) · ${pkg.installed_version || 'unknown'} → ${pkg.current_version}`,
          action: 'upgrade',
        })
      }
    }

    const actions: Array<{ name: string; description: string; action: string }> = [
      { name: 'Upgrade All', description: 'Upgrade all outdated packages', action: 'upgrade-all' },
      { name: 'Cleanup', description: 'Remove old versions and cache', action: 'cleanup' },
      { name: 'Doctor', description: 'Check system for problems', action: 'doctor' },
      { name: 'Refresh', description: 'Refresh installed and outdated lists', action: 'refresh' },
    ]

    for (const act of actions) {
      if (act.name.toLowerCase().includes(q)) {
        matches.push({ type: 'action', name: act.name, description: act.description, action: act.action })
      }
    }

    return matches.slice(0, 20)
  })

  return { query, isOpen, results, open, close }
}
