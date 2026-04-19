<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import ServicesPage from '@/pages/ServicesPage.vue'

const { t } = useI18n()
</script>

<template>
  <div class="flex h-full min-h-0 flex-col gap-4">
    <div class="content-header">
      <h1 class="content-title">{{ t('services.title') }}</h1>
    </div>
    <div class="min-h-0 flex-1">
      <ServicesPage embedded class="h-full" />
    </div>
  </div>
</template>
