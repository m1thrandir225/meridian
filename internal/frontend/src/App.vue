<script setup lang="ts">
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'
import { useAppearanceStore } from '@/stores/appearance'
import { Toaster } from './components/ui/sonner'
import 'vue-sonner/style.css'
import { VueQueryDevtools } from '@tanstack/vue-query-devtools'
const appearanceStore = useAppearanceStore()

onMounted(() => {
  // Initialize appearance settings on app mount
  appearanceStore.initializeAppearance()
})
</script>

<template>
  <RouterView />
  <Toaster />
  <VueQueryDevtools />
</template>
