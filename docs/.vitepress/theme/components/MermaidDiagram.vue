<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'
import { useData } from 'vitepress'

const props = defineProps<{
  code?: string
  encoded?: string
}>()

const { isDark } = useData()
const container = ref<HTMLElement | null>(null)
const errorMessage = ref('')
const diagramId = `mermaid-${Math.random().toString(36).slice(2)}`

let renderIndex = 0

async function renderDiagram() {
  if (!container.value) return

  errorMessage.value = ''
  const currentIndex = ++renderIndex

  try {
    const mermaid = (await import('mermaid')).default

    mermaid.initialize({
      startOnLoad: false,
      theme: isDark.value ? 'dark' : 'default',
      securityLevel: 'strict',
      fontFamily: 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif'
    })

    const id = `${diagramId}-${currentIndex}`
    const source = props.encoded ? decodeBase64Utf8(props.encoded) : props.code || ''
    const { svg } = await mermaid.render(id, source)

    if (currentIndex === renderIndex && container.value) {
      container.value.innerHTML = svg
    }
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Mermaid diagram render failed'
    if (container.value) {
      container.value.innerHTML = ''
    }
  }
}

function decodeBase64Utf8(value: string) {
  const binary = window.atob(value)
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0))

  return new TextDecoder().decode(bytes)
}

onMounted(renderDiagram)

watch(
  () => [props.code, props.encoded, isDark.value],
  async () => {
    await nextTick()
    await renderDiagram()
  }
)
</script>

<template>
  <figure class="mermaid-diagram">
    <div ref="container" class="mermaid-diagram__canvas" aria-label="技术图示"></div>
    <pre v-if="errorMessage" class="mermaid-diagram__error">{{ errorMessage }}</pre>
  </figure>
</template>
