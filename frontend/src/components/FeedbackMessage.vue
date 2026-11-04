<script setup lang="ts">
defineProps<{ feedback?: { text: string; error: boolean } }>()
</script>

<template>
  <p v-if="feedback" class="feedback-message" :class="{failed: feedback.error}" :role="feedback.error ? 'alert' : 'status'">{{feedback.text}}</p>
</template>

<style scoped>
.feedback-message { margin: .75rem 0 0; padding: .65rem .8rem; border-radius: 6px; color: #176b46; background: #edf5ef; overflow-wrap: anywhere }
.feedback-message.failed { color: #b42318; background: #fff0eb }
</style>
