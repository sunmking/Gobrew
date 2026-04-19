<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    rows?: number
    showHeader?: boolean
  }>(),
  {
    rows: 4,
    showHeader: true,
  },
)
</script>

<template>
  <div class="ui-skeleton-wrap">
    <div v-if="props.showHeader" class="ui-skeleton ui-skeleton-title" />
    <div class="ui-skeleton-list">
      <div v-for="i in props.rows" :key="i" class="ui-skeleton-row">
        <div class="ui-skeleton ui-skeleton-avatar" />
        <div class="ui-skeleton-content">
          <div class="ui-skeleton ui-skeleton-line w-40" />
          <div class="ui-skeleton ui-skeleton-line w-64" />
        </div>
      </div>
    </div>
  </div>
</template>
