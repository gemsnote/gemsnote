<script setup lang="ts">
import {ref, shallowRef, watch, nextTick, onMounted, onBeforeUnmount} from 'vue'
import {MermaidRenderer, previewDocument} from '../mermaidPreview'
import MermaidDiagram from './MermaidDiagram.vue'
const props = defineProps<{content: string; markdown: boolean}>()
const root = ref<HTMLElement>(), host = ref<HTMLElement>(), html = ref('')
const renderer = shallowRef<MermaidRenderer>()
const targets = shallowRef<{id: string; source: string; element: HTMLElement}[]>([])
let generation = 0
async function update() {
  const id = ++generation
  targets.value = []
  // Release old frames and pending work before replacing document content.
  renderer.value?.dispose()
  renderer.value = undefined
  const parsed = previewDocument(props.content, props.markdown, `mermaid-${Math.random().toString(36).slice(2)}`)
  html.value = parsed.html
  await nextTick()
  if (generation !== id || !root.value || !host.value) return
  if (parsed.diagrams.length) {
    renderer.value = new MermaidRenderer(host.value)
    targets.value = parsed.diagrams.map(d => ({...d, element: root.value!.querySelector<HTMLElement>(`#${d.id}`)!})).filter(d => d.element)
  }
}
onMounted(update)
watch(() => [props.content, props.markdown], update)
onBeforeUnmount(() => { generation++; renderer.value?.dispose() })
</script>
<template>
<div class="rendered">
 <div ref="root" v-html="html"></div>
 <div ref="host"></div>
 <template v-if="renderer"><Teleport v-for="target in targets" :key="target.id" :to="target.element"><MermaidDiagram :source="target.source" :renderer="renderer"/></Teleport></template>
</div>
</template>
