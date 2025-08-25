<template>
  <Card class="hover:shadow-lg transition-shadow duration-300">
    <CardHeader>
      <div class="text-4xl mb-2">{{ icon }}</div>
      <CardTitle>{{ title }}</CardTitle>
    </CardHeader>
    <CardContent>
      <CardDescription class="leading-relaxed">{{ description }}</CardDescription>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

defineProps<{
  icon: string
  title: string
  description: string
}>()
</script>
