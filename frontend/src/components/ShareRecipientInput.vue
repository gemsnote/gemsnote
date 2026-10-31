<script setup lang="ts">
import {computed,ref,useId,watch} from 'vue'
import {t} from '../i18n'
const props=defineProps<{modelValue:string;suggestions:string[];disabled?:boolean}>()
const emit=defineEmits<{'update:modelValue':[value:string]}>()
const focused=ref(false),dismissed=ref(false),active=ref(-1),listId=useId()
const candidates=computed(()=>[...new Set(props.suggestions.filter(Boolean))].filter(name=>name.toLowerCase().includes(props.modelValue.trim().toLowerCase())).sort((a,b)=>a.localeCompare(b)).slice(0,20))
const visible=computed(()=>focused.value&&!dismissed.value&&!props.disabled&&candidates.value.length>0)
watch(()=>props.modelValue,()=>{active.value=-1})
function select(name:string){emit('update:modelValue',name);dismissed.value=true;active.value=-1}
function keydown(event:KeyboardEvent){
 if(event.isComposing||event.keyCode===229)return
 if(event.key==='Escape'){dismissed.value=true;event.stopPropagation();return}
 if((event.key==='ArrowDown'||event.key==='ArrowUp')&&candidates.value.length){event.preventDefault();dismissed.value=false;active.value=event.key==='ArrowDown'?(active.value+1)%candidates.value.length:(active.value<0?candidates.value.length-1:(active.value+candidates.value.length-1)%candidates.value.length);return}
 if(event.key==='Enter'&&visible.value&&active.value>=0){event.preventDefault();select(candidates.value[active.value])}
}
</script>
<template>
<div class="share-recipient-input">
<input :value="modelValue" :disabled="disabled" required type="text" role="combobox" autocomplete="off" :aria-label="t('用户名或邮箱')" aria-autocomplete="list" :aria-expanded="visible" :aria-controls="listId" :aria-activedescendant="visible&&active>=0?`${listId}-${active}`:undefined" @input="emit('update:modelValue',($event.target as HTMLInputElement).value);dismissed=false" @focus="focused=true;dismissed=false" @blur="focused=false" @keydown="keydown">
<ul v-if="visible" :id="listId" role="listbox" :aria-label="t('已共享过的用户')">
<li v-for="(name,index) in candidates" :id="`${listId}-${index}`" :key="name" role="option" :aria-selected="active===index" :class="{active:active===index}" @pointerdown.prevent @click="select(name)">{{name}}</li>
</ul>
</div>
</template>
<style scoped>
.share-recipient-input{position:relative}.share-recipient-input input{width:100%;box-sizing:border-box}.share-recipient-input ul{position:absolute;top:100%;left:0;right:0;z-index:20;max-height:220px;overflow:auto;list-style:none;padding:4px;margin:3px 0;background:white;border:1px solid #d9e0da;border-radius:8px;box-shadow:0 4px 16px #193b3526}.share-recipient-input li{padding:8px;cursor:pointer;overflow-wrap:anywhere;border-radius:5px}.share-recipient-input li:hover,.share-recipient-input li.active{background:#dbe9dd;color:#193b35}
</style>
