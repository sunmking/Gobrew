<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Play, RefreshCw, RotateCw, Square } from 'lucide-vue-next'
import BrewButton from '@/components/common/BrewButton.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import { useServicesStore } from '@/stores/services'

const servicesStore = useServicesStore()
const { t } = useI18n()

function statusKind(status: string) {
  const lower = status.toLowerCase()
  if (lower.includes('started') || lower.includes('running')) return 'running'
  if (lower.includes('error')) return 'error'
  return 'stopped'
}

onMounted(servicesStore.fetchServices)
</script>

<template>
  <section class="page">
    <div class="content-header">
      <h1 class="content-title">{{ t('services.title') }}</h1>
      <p class="content-subtitle">{{ t('servicesPage.subtitle') }}</p>
    </div>
    <div class="toolbar">
      <BrewButton variant="primary" @click="servicesStore.restartAll"><RotateCw :size="14" />{{ t('services.restartAll') }}</BrewButton>
      <BrewButton @click="servicesStore.fetchServices"><RefreshCw :size="14" />{{ t('common.refresh') }}</BrewButton>
    </div>
    <div class="content-body" style="padding-top:0;">
      <table class="pkg-table">
        <thead><tr><th>{{ t('servicesPage.name') }}</th><th>{{ t('table.status') }}</th><th>{{ t('servicesPage.user') }}</th><th>{{ t('servicesPage.file') }}</th><th>{{ t('servicesPage.exitCode') }}</th><th></th></tr></thead>
        <tbody>
          <tr v-for="service in servicesStore.services" :key="service.name">
            <td class="pkg-name">{{ service.name }}</td>
            <td><StatusPill :status="statusKind(service.status)">{{ service.status }}</StatusPill></td>
            <td>{{ service.user || '-' }}</td>
            <td class="pkg-desc">{{ service.file || '-' }}</td>
            <td class="pkg-version">{{ service.exit_code }}</td>
            <td style="text-align:right;">
              <div style="display:flex; gap:6px; justify-content:flex-end;">
                <BrewButton variant="ghost" @click="servicesStore.start(service.name)"><Play :size="13" />{{ t('services.start') }}</BrewButton>
                <BrewButton variant="ghost" @click="servicesStore.stop(service.name)"><Square :size="13" />{{ t('services.stop') }}</BrewButton>
                <BrewButton variant="ghost" @click="servicesStore.restart(service.name)"><RotateCw :size="13" />{{ t('services.restart') }}</BrewButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="servicesStore.services.length === 0" class="empty-state">{{ t('services.noServices') }}</div>
    </div>
  </section>
</template>
