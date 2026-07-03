<script setup lang="ts">
interface LearningStep {
  title: string
  description: string
  link?: string
  badge?: string
}

defineProps<{
  steps: LearningStep[]
}>()
</script>

<template>
  <ol class="learning-path">
    <li v-for="(step, index) in steps" :key="step.title" class="learning-path__item">
      <a v-if="step.link" class="learning-path__card" :href="step.link">
        <span class="learning-path__index">{{ String(index + 1).padStart(2, '0') }}</span>
        <span class="learning-path__content">
          <span class="learning-path__title">{{ step.title }}</span>
          <span class="learning-path__description">{{ step.description }}</span>
        </span>
        <span v-if="step.badge" class="learning-path__badge">{{ step.badge }}</span>
      </a>
      <div v-else class="learning-path__card">
        <span class="learning-path__index">{{ String(index + 1).padStart(2, '0') }}</span>
        <span class="learning-path__content">
          <span class="learning-path__title">{{ step.title }}</span>
          <span class="learning-path__description">{{ step.description }}</span>
        </span>
        <span v-if="step.badge" class="learning-path__badge">{{ step.badge }}</span>
      </div>
    </li>
  </ol>
</template>
