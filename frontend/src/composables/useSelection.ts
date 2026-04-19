import { computed, ref } from 'vue'

export function useSelection() {
  const selected = ref<string[]>([])

  const selectedCount = computed(() => selected.value.length)

  function isSelected(id: string) {
    return selected.value.includes(id)
  }

  function toggle(id: string) {
    if (isSelected(id)) {
      selected.value = selected.value.filter((item) => item !== id)
      return
    }
    selected.value = [...selected.value, id]
  }

  function clear() {
    selected.value = []
  }

  function selectAll(ids: string[]) {
    selected.value = [...ids]
  }

  function isAllSelected(ids: string[]) {
    return ids.length > 0 && ids.every((id) => selected.value.includes(id))
  }

  function sync(ids: string[]) {
    const allowed = new Set(ids)
    selected.value = selected.value.filter((id) => allowed.has(id))
  }

  return {
    selected,
    selectedCount,
    isSelected,
    toggle,
    clear,
    selectAll,
    isAllSelected,
    sync,
  }
}
