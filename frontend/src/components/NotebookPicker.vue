<script setup lang="ts">
import {computed,ref} from 'vue'
import {t} from '../i18n'
const props=defineProps<{items:any[];modelValue:string;exclude?:string;root?:boolean}>()
const emit=defineEmits<{ 'update:modelValue':[value:string] }>()
const search=ref(''),collapsed=ref(new Set<string>())
const rows=computed(()=>{
 const key=search.value.trim().toLowerCase()
 function visit(items:any[],depth=0):any[]{return items.flatMap(n=>{
  if(n.NotebookId===props.exclude)return []
  const children=visit(n.Subs||[],depth+1),match=String(n.Title).toLowerCase().includes(key)
  if(key&&!match&&!children.length)return []
  return [{...n,depth},...(key||!collapsed.value.has(n.NotebookId)?children:[])]
 })}
 return visit(props.items)
})
function toggle(id:string){const next=new Set(collapsed.value);next.has(id)?next.delete(id):next.add(id);collapsed.value=next}
</script>
<template>
<div class="notebook-picker">
<input v-model="search" :placeholder="t('搜索笔记本')" :aria-label="t('搜索笔记本')">
<div class="picker-tree" :aria-label="t('选择笔记本')">
<button v-if="root" type="button" :class="{selected:modelValue===''}" @click="emit('update:modelValue','')">{{t('根目录')}}</button>
<div v-for="n in rows" :key="n.NotebookId" class="picker-row" :style="{paddingLeft:n.depth*18+'px'}">
<button v-if="n.Subs?.length" type="button" :aria-label="t('展开或收起笔记本')" :aria-expanded="!collapsed.has(n.NotebookId)" @click="toggle(n.NotebookId)">{{collapsed.has(n.NotebookId)?'▸':'▾'}}</button><span v-else class="picker-leaf">·</span>
<button type="button" :class="{selected:modelValue===n.NotebookId}" :aria-pressed="modelValue===n.NotebookId" :title="n.Title" @click="emit('update:modelValue',n.NotebookId)">{{n.Title}}</button>
</div>
<p v-if="!rows.length">{{t('没有匹配的笔记本')}}</p>
</div></div>
</template>
<style scoped>
.notebook-picker input{width:100%;box-sizing:border-box}.picker-tree{max-height:300px;overflow:auto;margin:10px 0;border:1px solid #d9e0da;border-radius:7px;padding:5px}.picker-row{display:flex;align-items:center}.picker-row button:last-child{flex:1;text-align:left;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.picker-leaf{width:32px;text-align:center;flex-shrink:0}.picker-row button{border:0;background:transparent}.picker-tree>button{width:100%;text-align:left}
</style>
