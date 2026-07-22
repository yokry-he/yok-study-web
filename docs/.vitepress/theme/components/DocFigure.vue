<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  src: string
  alt: string
  caption: string
  sourceUrl?: string
  sourceLabel?: string
  width?: number
  height?: number
  zoomable?: boolean
}>(), {
  sourceLabel: '图片来源',
  zoomable: true
})

const open = ref(false)
const failed = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const dialog = ref<HTMLElement | null>(null)
const closeButton = ref<HTMLButtonElement | null>(null)
let previousBodyOverflow = ''

function close() {
  open.value = false
}

function onDialogKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    event.preventDefault()
    close()
    return
  }
  if (event.key !== 'Tab' || !dialog.value) return

  const focusable = [...dialog.value.querySelectorAll<HTMLElement>(
    'button, [href], [tabindex]:not([tabindex="-1"])'
  )]
  if (focusable.length === 0) return

  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

watch(open, async value => {
  if (typeof document === 'undefined') return
  if (value) {
    previousBodyOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    await nextTick()
    closeButton.value?.focus()
    return
  }

  document.body.style.overflow = previousBodyOverflow
  trigger.value?.focus()
})

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = previousBodyOverflow
  }
})
</script>

<template>
  <figure class="doc-figure">
    <button
      v-if="props.zoomable && !failed"
      ref="trigger"
      class="doc-figure__trigger"
      type="button"
      :aria-label="`放大查看：${props.alt}`"
      @click="open = true"
    >
      <img
        class="doc-figure__image"
        :src="props.src"
        :alt="props.alt"
        :width="props.width"
        :height="props.height"
        loading="lazy"
        decoding="async"
        @error="failed = true"
      >
    </button>
    <img
      v-else-if="!failed"
      class="doc-figure__image"
      :src="props.src"
      :alt="props.alt"
      :width="props.width"
      :height="props.height"
      loading="lazy"
      decoding="async"
      @error="failed = true"
    >
    <p v-else class="doc-figure__error" role="status">
      图片加载失败，请根据图注继续阅读。
    </p>
    <figcaption class="doc-figure__caption">
      {{ props.caption }}
      <a
        v-if="props.sourceUrl"
        class="doc-figure__source"
        :href="props.sourceUrl"
        target="_blank"
        rel="noreferrer"
      >{{ props.sourceLabel }}</a>
    </figcaption>

    <Teleport to="body">
      <div
        v-if="open"
        ref="dialog"
        class="doc-figure-lightbox"
        role="dialog"
        aria-modal="true"
        :aria-label="props.alt"
        @click.self="close"
        @keydown="onDialogKeydown"
      >
        <button
          ref="closeButton"
          class="doc-figure-lightbox__close"
          type="button"
          aria-label="关闭图片预览"
          @click="close"
        >
          ×
        </button>
        <img class="doc-figure-lightbox__image" :src="props.src" :alt="props.alt">
      </div>
    </Teleport>
  </figure>
</template>
