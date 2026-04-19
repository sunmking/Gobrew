<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  description?: string
  badge?: string
  badgeColor?: 'blue' | 'green' | 'amber' | 'red'
}>()
</script>

<template>
  <div class="pkg-card rounded-md p-3 transition-colors">
    <div class="flex items-center justify-between gap-2">
      <div class="min-w-0">
        <div class="flex items-center gap-1.5">
          <h3 class="pkg-title truncate text-sm font-medium">{{ title }}</h3>
          <span
            v-if="badge"
            class="pkg-badge inline-flex flex-shrink-0 items-center rounded-full px-1.5 py-0.5 text-[10px] font-medium leading-none"
            :class="{
              'pkg-badge--blue': badgeColor === 'blue' || !badgeColor,
              'pkg-badge--green': badgeColor === 'green',
              'pkg-badge--amber': badgeColor === 'amber',
              'pkg-badge--red': badgeColor === 'red',
            }"
          >{{ badge }}</span>
        </div>
        <p v-if="subtitle" class="pkg-subtitle truncate text-xs">{{ subtitle }}</p>
        <p v-if="description" class="pkg-description mt-1 line-clamp-2 text-xs">{{ description }}</p>
      </div>
      <div class="flex items-center gap-1 flex-shrink-0 ml-2">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.pkg-card {
  background: var(--color-card);
  border: 1px solid var(--color-border);
}

.pkg-card:hover {
  border-color: var(--color-accent);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.pkg-title {
  color: var(--color-text);
}

.pkg-subtitle {
  color: var(--color-text-secondary);
}

.pkg-description {
  color: var(--color-text-tertiary);
}

.pkg-badge {
  border: 1px solid;
}

.pkg-badge--blue {
  background: var(--color-accent-light);
  color: var(--color-accent);
  border-color: var(--color-accent);
}

.pkg-badge--green {
  background: var(--color-success-light);
  color: var(--color-success);
  border-color: var(--color-success);
}

.pkg-badge--amber {
  background: var(--color-warning-light);
  color: var(--color-warning);
  border-color: var(--color-warning);
}

.pkg-badge--red {
  background: var(--color-danger-light);
  color: var(--color-danger);
  border-color: var(--color-danger);
}
</style>
