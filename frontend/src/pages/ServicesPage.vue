<script setup lang="ts">
import { onMounted } from 'vue'
import { Play, RefreshCw, RotateCw, Square } from 'lucide-vue-next'
import BrewButton from '@/components/common/BrewButton.vue'
import StatusPill from '@/components/common/StatusPill.vue'
import { useServicesStore } from '@/stores/services'

const servicesStore = useServicesStore()

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
      <h1 class="content-title">Services</h1>
      <p class="content-subtitle">管理 Homebrew 后台服务的启动、停止和重启</p>
    </div>
    <div class="toolbar">
      <BrewButton variant="primary" @click="servicesStore.restartAll"><RotateCw :size="14" />全部重启</BrewButton>
      <BrewButton @click="servicesStore.fetchServices"><RefreshCw :size="14" />刷新状态</BrewButton>
    </div>
    <div class="content-body" style="padding-top:0;">
      <table class="pkg-table">
        <thead><tr><th>服务名称</th><th>状态</th><th>用户</th><th>文件</th><th>退出码</th><th></th></tr></thead>
        <tbody>
          <tr v-for="service in servicesStore.services" :key="service.name">
            <td class="pkg-name">{{ service.name }}</td>
            <td><StatusPill :status="statusKind(service.status)">{{ service.status }}</StatusPill></td>
            <td>{{ service.user || '-' }}</td>
            <td class="pkg-desc">{{ service.file || '-' }}</td>
            <td class="pkg-version">{{ service.exit_code }}</td>
            <td style="text-align:right;">
              <div style="display:flex; gap:6px; justify-content:flex-end;">
                <BrewButton variant="ghost" @click="servicesStore.start(service.name)"><Play :size="13" />启动</BrewButton>
                <BrewButton variant="ghost" @click="servicesStore.stop(service.name)"><Square :size="13" />停止</BrewButton>
                <BrewButton variant="ghost" @click="servicesStore.restart(service.name)"><RotateCw :size="13" />重启</BrewButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="servicesStore.services.length === 0" class="empty-state">暂无 Homebrew services</div>
    </div>
  </section>
</template>
