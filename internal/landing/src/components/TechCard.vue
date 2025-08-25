<template>
  <Card class="text-center group hover:shadow-lg transition-all duration-300">
    <CardContent class="pt-6">
      <div
        class="w-16 h-16 mx-auto mb-4 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-lg group-hover:scale-110 transition-transform duration-300"
        :class="color"
      >
        <img class="object-contain w-10 h-10" :src="`/${name.toLowerCase()}.svg`" />
      </div>
      <CardTitle class="text-lg">{{ name }}</CardTitle>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { Card, CardContent, CardTitle } from '@/components/ui/card'

defineProps<{
  name: string
  color: string
}>()
</script>
