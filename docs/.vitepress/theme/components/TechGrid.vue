<script setup lang="ts">
interface TechItem {
  title: string
  description: string
  link: string
  level?: string
}

defineProps<{
  items: TechItem[]
}>()
</script>

<template>
  <div class="tech-grid">
    <a v-for="item in items" :key="item.title" class="tech-card" :href="item.link">
      <span v-if="item.level" class="tech-card__level">{{ item.level }}</span>
      <strong class="tech-card__title">{{ item.title }}</strong>
      <span class="tech-card__description">{{ item.description }}</span>
    </a>
  </div>
</template>
