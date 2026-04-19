<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    label?: string
    compact?: boolean
  }>(),
  {
    label: '',
    compact: false,
  },
)

const { t } = useI18n()
</script>

<template>
  <div class="ui-loading-inline" :class="{ 'is-compact': props.compact }">
    <Loader2 :size="14" class="ui-spinner" />
    <span>{{ props.label || t('common.loading') }}</span>
  </div>
</template>
