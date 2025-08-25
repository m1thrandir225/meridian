<script setup lang="ts">
import { useAppearanceStore } from '@/stores/appearance'
import type { Theme, AccentColor } from '@/stores/appearance'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Palette, Sun, Moon, Monitor } from 'lucide-vue-next'

const appearanceStore = useAppearanceStore()

const themes: { value: Theme; label: string; icon: any }[] = [
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
  { value: 'system', label: 'System', icon: Monitor },
]

const accentColors: { value: AccentColor; label: string; color: string }[] = [
  { value: 'blue', label: 'Blue', color: 'bg-blue-500' },
  { value: 'green', label: 'Green', color: 'bg-green-500' },
  { value: 'purple', label: 'Purple', color: 'bg-purple-500' },
  { value: 'red', label: 'Red', color: 'bg-red-500' },
  { value: 'orange', label: 'Orange', color: 'bg-orange-500' },
]
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon" class="h-9 w-9">
        <Palette class="h-4 w-4" />
        <span class="sr-only">Toggle appearance</span>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" class="w-56">
      <DropdownMenuLabel>Theme</DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuItem
        v-for="theme in themes"
        :key="theme.value"
        @click="appearanceStore.setTheme(theme.value)"
        :class="['flex items-center gap-2', appearanceStore.theme === theme.value && 'bg-accent']"
      >
        <component :is="theme.icon" class="h-4 w-4" />
        {{ theme.label }}
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuLabel>Accent Color</DropdownMenuLabel>
      <DropdownMenuSeparator />
      <div class="flex gap-2 p-2">
        <button
          v-for="color in accentColors"
          :key="color.value"
          @click="appearanceStore.setAccentColor(color.value)"
          :class="[
            'w-8 h-8 rounded-full border-2 transition-all',
            color.color,
            appearanceStore.accentColor === color.value
              ? 'ring-2 ring-offset-2 ring-primary scale-110'
              : 'hover:scale-105',
          ]"
          :title="color.label"
        >
          <div
            v-if="appearanceStore.accentColor === color.value"
            class="w-full h-full flex items-center justify-center"
          >
            <div class="w-2 h-2 bg-white rounded-full" />
          </div>
        </button>
      </div>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
