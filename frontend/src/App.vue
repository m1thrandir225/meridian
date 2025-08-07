<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAppearanceStore } from '@/stores/appearance'

const appearanceStore = useAppearanceStore()

onMounted(() => {
  // Initialize appearance settings on app mount
  appearanceStore.initializeAppearance()
})
</script>

<template>
  <RouterView />
</template>
