<script setup lang="ts">
import {ref, watch, onBeforeUnmount} from 'vue'
import {type DiagramImage, type MermaidRenderer} from '../mermaidPreview'
import {t} from '../i18n'
const props = defineProps<{source: string; renderer: MermaidRenderer}>()
const result = ref<DiagramImage>(), failed = ref(false), expanded = ref(false), scale = ref(100)
let generation = 0
watch(() => props.source, async source => {
  const id = ++generation
  result.value = undefined; failed.value = false
  try { const image = await props.renderer.render(source); if (id === generation) result.value = image }
  catch { if (id === generation) failed.value = true }
}, {immediate: true})
onBeforeUnmount(() => { generation++ })
function imageURL() { return 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(result.value!.svg) }
</script>
<template>
<figure class="mermaid-diagram" :aria-label="t('Mermaid 图表')">
  <p v-if="failed" role="status">{{t('图表无法渲染，请检查语法或图表大小。')}}</p>
  <p v-else-if="!result" role="status">{{t('正在绘制图表…')}}</p>
  <button v-if="result" class="diagram-image" :aria-label="t('放大图表')" @click="expanded=true;scale=100"><img :src="imageURL()" :alt="t('Mermaid 图表')" :style="{width:result.width+'px'}"></button>
  <details :open="failed"><summary>{{t('图表源码')}}</summary><pre><code>{{source}}</code></pre></details>
</figure>
<Teleport to="body"><div v-if="expanded&&result" class="diagram-backdrop" @click.self="expanded=false" @keydown.esc="expanded=false">
 <section role="dialog" aria-modal="true" :aria-label="t('Mermaid 图表')" class="diagram-dialog">
  <header><label>{{t('缩放')}} <input v-model="scale" type="range" min="25" max="300" step="25"></label><button autofocus @click="expanded=false">{{t('关闭')}}</button></header>
  <div class="diagram-scroll"><img :src="imageURL()" :alt="t('Mermaid 图表')" :style="{width:result.width*scale/100+'px'}"></div>
 </section>
</div></Teleport>
</template>
<style scoped>
.mermaid-diagram{margin:1em 0;padding:12px;border:1px solid #d5ded9;border-radius:8px;background:#fff;color:#203d34;max-width:100%;overflow:hidden}
.diagram-image{display:block;background:transparent;border:0;padding:0;max-width:100%;cursor:zoom-in}.diagram-image img{display:block;max-width:100%;height:auto}
summary{cursor:pointer;font-size:.85em;margin-top:8px}pre{overflow:auto;white-space:pre;font-size:.85em}p{font-size:.9em}
.diagram-backdrop{position:fixed;inset:0;z-index:1500;background:#0008;display:grid;place-items:center}.diagram-dialog{width:94vw;height:90vh;display:flex;flex-direction:column;background:white;border-radius:8px;padding:16px}.diagram-dialog header{display:flex;justify-content:space-between;align-items:center;gap:16px}.diagram-dialog label{display:flex;align-items:center;gap:8px}.diagram-scroll{overflow:auto;flex:1;margin-top:12px}.diagram-scroll img{max-width:none}
</style>
