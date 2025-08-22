<template>
  <div class="text-center group">
    <div
      class="w-16 h-16 mx-auto mb-4 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-lg group-hover:scale-110 transition-transform duration-300"
      :class="color"
    >
      <img class="object-contain w-10 h-10" :src="`/${name.toLowerCase()}.svg`" />
    </div>
    <h3 class="text-lg font-semibold text-slate-900">{{ name }}</h3>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  name: string
  color: string
}>()
</script>
