<script setup lang="ts">
import {computed,nextTick,ref,useId,watch} from 'vue'
import {language,t} from '../i18n'
const props=defineProps<{modelValue:string; suggestions:string[]; readonly?:boolean}>()
const emit=defineEmits<{'update:modelValue':[value:string]}>()
const draft=ref(''),focused=ref(false),dismissed=ref(false),active=ref(-1),input=ref<HTMLInputElement>()
const listId=useId()
let composing=false
function compositionStart(){composing=true}
function compositionEnd(){composing=false;typed()}
const tokens=computed(()=>[...new Set(props.modelValue.split(',').map(s=>s.trim()).filter(Boolean))])
const candidates=computed(()=>[...new Set(props.suggestions.map(s=>s.trim()).filter(s=>s&&!tokens.value.includes(s)&&!s.includes(',')))]
 .filter(s=>s.toLocaleLowerCase().includes(draft.value.trim().toLocaleLowerCase()))
 .sort((a,b)=>a.localeCompare(b,language.value,{numeric:true,sensitivity:'base'})).slice(0,20))
const visible=computed(()=>focused.value&&!props.readonly&&!dismissed.value&&candidates.value.length>0)
watch([draft,candidates],()=>{active.value=-1})
function publish(values:string[]){const value=[...new Set(values.map(s=>s.trim()).filter(Boolean))].join(',');if(value!==props.modelValue)emit('update:modelValue',value)}
function commit(value=draft.value){if(props.readonly||composing)return;publish([...tokens.value,...value.split(/[,，\n]/)]);draft.value='';active.value=-1;dismissed.value=false}
function choose(value:string){commit(value);input.value?.focus()}
function remove(value:string){if(!props.readonly)publish(tokens.value.filter(s=>s!==value))}
function typed(){dismissed.value=false;if(composing)return;const parts=draft.value.split(/[,，\n]/);if(parts.length>1){const tail=parts.pop()!;publish([...tokens.value,...parts]);draft.value=tail}}
function blur(){focused.value=false;commit()}
function keydown(event:KeyboardEvent){
 if(props.readonly||composing||event.isComposing||event.keyCode===229)return
 if(event.key==='Escape'){dismissed.value=true;active.value=-1;event.stopPropagation();return}
 if((event.key==='ArrowDown'||event.key==='ArrowUp')&&candidates.value.length){
  event.preventDefault();dismissed.value=false
  active.value=event.key==='ArrowDown'?(active.value+1)%candidates.value.length:(active.value<0?candidates.value.length-1:(active.value+candidates.value.length-1)%candidates.value.length)
  nextTick(()=>document.getElementById(`${listId}-${active.value}`)?.scrollIntoView({block:'nearest'}));return
 }
 if(event.key==='Enter'){event.preventDefault();commit(visible.value&&active.value>=0?candidates.value[active.value]:draft.value)}
 if(event.key==='Backspace'&&!draft.value&&tokens.value.length){event.preventDefault();remove(tokens.value.at(-1)!)}
}
defineExpose({commit})
</script>
<template>
<div class="tag-composer">
<div class="tag-chips">
<span v-for="tag in tokens" :key="tag" class="tag-chip"><span :title="tag">{{tag}}</span><button v-if="!readonly" type="button" :aria-label="t('删除标签 {tag}',{tag})" @pointerdown.prevent @click="remove(tag)">×</button></span>
<input v-if="!readonly" ref="input" v-model="draft" class="tag-input" role="combobox" aria-autocomplete="list" :aria-label="t('笔记标签')" :aria-controls="listId" :aria-expanded="visible" :aria-activedescendant="visible&&active>=0?`${listId}-${active}`:undefined" :placeholder="t('输入标签，回车确认')" autocomplete="off" @input="typed" @keydown="keydown" @focus="focused=true;dismissed=false" @blur="blur" @compositionstart="compositionStart" @compositionend="compositionEnd">
</div>
<ul v-if="visible" :id="listId" role="listbox" :aria-label="t('已有标签')" class="tag-candidates">
<li v-for="(tag,index) in candidates" :id="`${listId}-${index}`" :key="tag" role="option" :aria-selected="active===index" :class="{active:active===index}" @pointerdown.prevent @click="choose(tag)">{{tag}}</li>
</ul>
</div>
</template>
<style scoped>
.tag-composer{position:relative;flex:1;min-width:0}.tag-chips{display:flex;gap:5px;align-items:center;flex-wrap:wrap;max-height:120px;overflow:auto;padding:2px}.tag-chip{display:inline-flex;align-items:center;gap:3px;max-width:100%;border-radius:6px;background:#dbe9dd;color:#193b35;padding:3px 7px}.tag-chip>span{overflow:hidden;text-overflow:ellipsis;white-space:nowrap;min-width:0;flex:1}.tag-chip button{border:0;background:transparent;padding:0 3px;color:inherit;line-height:1.2}.tag-input{min-width:140px;flex:1}.tag-candidates{position:absolute;bottom:calc(100% + 6px);left:0;z-index:12;width:min(320px,100%);max-height:220px;overflow:auto;list-style:none;margin:0;padding:4px;background:white;border:1px solid #d9e0da;border-radius:8px;box-shadow:0 4px 16px #193b3526}.tag-candidates li{padding:7px 10px;border-radius:5px;cursor:pointer;overflow-wrap:anywhere}.tag-candidates li:hover,.tag-candidates li.active{background:#dbe9dd;color:#193b35}
</style>
