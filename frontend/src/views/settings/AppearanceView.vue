<script setup lang="ts">
import SettingsLayout from '@/layouts/SettingsLayout.vue'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useAppearanceStore } from '@/stores/appearance'
import type { Theme, AccentColor, FontSize, MessageDisplayMode } from '@/stores/appearance'
import { ref } from 'vue'
import { toast } from 'vue-sonner'

const scrollContainer = ref<HTMLDivElement | null>(null)

const appearanceStore = useAppearanceStore()

const themes: { value: Theme; label: string }[] = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'system', label: 'System' },
]

const accentColors: { value: AccentColor; label: string; color: string }[] = [
  { value: 'blue', label: 'Blue', color: 'bg-blue-500' },
  { value: 'green', label: 'Green', color: 'bg-green-500' },
  { value: 'purple', label: 'Purple', color: 'bg-purple-500' },
  { value: 'red', label: 'Red', color: 'bg-red-500' },
  { value: 'orange', label: 'Orange', color: 'bg-orange-500' },
]

const fontSizes: { value: FontSize; label: string }[] = [
  { value: 'small', label: 'Small' },
  { value: 'medium', label: 'Medium' },
  { value: 'large', label: 'Large' },
]

const messageDisplayModes: { value: MessageDisplayMode; label: string; description: string }[] = [
  { value: 'cozy', label: 'Cozy', description: 'Modern and comfortable spacing' },
  { value: 'compact', label: 'Compact', description: 'Fit more messages on screen' },
]

const scrollToTopContainer = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }
}

const handleReset = () => {
  appearanceStore.resetToDefaults()
  scrollToTopContainer()
  toast.success('Appearance settings reset to defaults')
}

const onSave = () => {
  toast.success('Appearance settings saved')
  scrollToTopContainer()
}
</script>

<template>
  <SettingsLayout>
    <div class="flex flex-col h-full">
      <!-- Header -->
      <header
        class="flex items-center justify-between px-6 py-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
      >
        <div>
          <h1 class="text-2xl font-semibold">Appearance</h1>
          <p class="text-sm text-muted-foreground">Customize your chat experience</p>
        </div>
      </header>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6" ref="scrollContainer">
        <div class="max-w-2xl mx-auto space-y-6">
          <!-- Theme -->
          <Card>
            <CardHeader>
              <CardTitle>Theme</CardTitle>
              <CardDescription>Choose your preferred color scheme</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="grid grid-cols-3 gap-3">
                <div
                  v-for="theme in themes"
                  :key="theme.value"
                  :class="[
                    'relative border rounded-lg p-3 cursor-pointer hover:bg-accent transition-colors',
                    appearanceStore.theme === theme.value && 'ring-2 ring-primary',
                  ]"
                  @click="appearanceStore.setTheme(theme.value)"
                >
                  <div class="flex items-center space-x-2">
                    <div
                      :class="[
                        'w-4 h-4 rounded-full border-2',
                        appearanceStore.theme === theme.value
                          ? 'border-primary bg-primary'
                          : 'border-muted-foreground',
                      ]"
                    />
                    <Label class="cursor-pointer">{{ theme.label }}</Label>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Accent Color -->
          <Card>
            <CardHeader>
              <CardTitle>Accent Color</CardTitle>
              <CardDescription>Choose your preferred accent color</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="flex gap-3">
                <div
                  v-for="color in accentColors"
                  :key="color.value"
                  :class="[
                    'relative w-12 h-12 rounded-lg cursor-pointer border-2 transition-all',
                    color.color,
                    appearanceStore.accentColor === color.value
                      ? 'ring-2 ring-offset-2 ring-primary scale-110'
                      : 'hover:scale-105',
                  ]"
                  :title="color.label"
                  @click="appearanceStore.setAccentColor(color.value)"
                >
                  <div
                    v-if="appearanceStore.accentColor === color.value"
                    class="absolute inset-0 flex items-center justify-center"
                  >
                    <div class="w-4 h-4 bg-white rounded-full" />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Font Size -->
          <Card>
            <CardHeader>
              <CardTitle>Font Size</CardTitle>
              <CardDescription>Adjust text size for better readability</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="space-y-3">
                <div
                  v-for="size in fontSizes"
                  :key="size.value"
                  :class="[
                    'flex items-center space-x-3 border rounded-lg p-3 cursor-pointer hover:bg-accent transition-colors',
                    appearanceStore.fontSize === size.value && 'ring-2 ring-primary',
                  ]"
                  @click="appearanceStore.setFontSize(size.value)"
                >
                  <div
                    :class="[
                      'w-4 h-4 rounded-full border-2',
                      appearanceStore.fontSize === size.value
                        ? 'border-primary bg-primary'
                        : 'border-muted-foreground',
                    ]"
                  />
                  <Label class="cursor-pointer">{{ size.label }}</Label>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Message Display -->
          <Card>
            <CardHeader>
              <CardTitle>Message Display</CardTitle>
              <CardDescription>Choose how messages are displayed</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="space-y-3">
                <div
                  v-for="mode in messageDisplayModes"
                  :key="mode.value"
                  :class="[
                    'border rounded-lg p-4 cursor-pointer hover:bg-accent transition-colors',
                    appearanceStore.messageDisplayMode === mode.value && 'ring-2 ring-primary',
                  ]"
                  @click="appearanceStore.setMessageDisplayMode(mode.value)"
                >
                  <div class="flex items-start space-x-3">
                    <div
                      :class="[
                        'w-4 h-4 rounded-full border-2 mt-0.5',
                        appearanceStore.messageDisplayMode === mode.value
                          ? 'border-primary bg-primary'
                          : 'border-muted-foreground',
                      ]"
                    />
                    <div>
                      <Label class="cursor-pointer font-medium">{{ mode.label }}</Label>
                      <p class="text-sm text-muted-foreground mt-1">{{ mode.description }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- Actions -->
          <div class="flex justify-end gap-2">
            <Button variant="outline" @click="handleReset">Reset to Default</Button>
            <Button @click="onSave">Save Changes</Button>
          </div>
        </div>
      </div>
    </div>
  </SettingsLayout>
</template>
