<script setup lang="ts">
import {computed,onBeforeUnmount,onMounted,ref,watch} from 'vue'
const props=defineProps<{items:any[]}>()
const root=ref<HTMLElement>(),top=ref(0),height=ref(600)
const rowHeight=112,overscan=6
const virtual=computed(()=>props.items.length>300)
const start=computed(()=>virtual.value?Math.max(0,Math.floor(top.value/rowHeight)-overscan):0)
const end=computed(()=>virtual.value?Math.min(props.items.length,start.value+Math.ceil(height.value/rowHeight)+overscan*2):props.items.length)
const visible=computed(()=>props.items.slice(start.value,end.value))
let observer:ResizeObserver|undefined
onMounted(()=>{observer=new ResizeObserver(()=>height.value=root.value?.clientHeight||600);if(root.value)observer.observe(root.value)})
onBeforeUnmount(()=>observer?.disconnect())
watch([()=>props.items[0]?.NoteId,()=>props.items.length],()=>{top.value=0;if(root.value)root.value.scrollTop=0})
</script>
<template>
<div ref="root" class="list-items" :style="{'--note-row-height':rowHeight+'px'}" :class="{'virtual-note-list':virtual}" @scroll.passive="top=($event.target as HTMLElement).scrollTop">
<div v-if="virtual" aria-hidden="true" :style="{height:start*rowHeight+'px'}"></div>
<template v-for="item in visible" :key="item.NoteId"><slot :item="item"/></template>
<div v-if="virtual" aria-hidden="true" :style="{height:(items.length-end)*rowHeight+'px'}"></div>
<slot name="empty"/>
</div>
</template>
<style>
.virtual-note-list .note-row{height:var(--note-row-height);box-sizing:border-box}.virtual-note-list .note-item{height:100%;box-sizing:border-box}.virtual-note-list .note-item p{margin:6px 0}.virtual-note-list .note-item small{display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
</style>
