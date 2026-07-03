<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  tone?: 'mint' | 'sky' | 'lemon' | 'rose'
}>(), {
  tone: 'mint'
})
</script>

<template>
  <section class="practice-block" :class="`practice-block--${tone}`">
    <div class="practice-block__label">实践提示</div>
    <h3 class="practice-block__title">{{ title }}</h3>
    <div class="practice-block__body">
      <slot />
    </div>
  </section>
</template>
