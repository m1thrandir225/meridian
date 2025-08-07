<script setup lang="ts">
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Settings, User, Lock, Palette, ArrowLeft } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
import { ref } from 'vue'

import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarFooter,
  SidebarRail,
} from '@/components/ui/sidebar'

const isSettingsSidebarOpen = ref(true)

// Sample user data
const user = {
  name: 'John Doe',
  username: 'johndoe',
  avatar: '/avatars/user.png',
}

const settingsNavigation = [
  {
    name: 'Profile',
    href: '/settings/profile',
    icon: User,
  },
  {
    name: 'Password',
    href: '/settings/password',
    icon: Lock,
  },
  {
    name: 'Appearance',
    href: '/settings/appearance',
    icon: Palette,
  },
]
</script>

<template>
  <div class="flex h-screen">
    <SidebarProvider>
      <Transition
        name="slide-left"
        enter-active-class="transition-all duration-300 ease-out"
        leave-active-class="transition-all duration-300 ease-in"
        enter-from-class="-translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="-translate-x-full opacity-0"
      >
        <Sidebar v-if="isSettingsSidebarOpen" class="border-r">
          <SidebarHeader>
            <div class="flex items-center gap-3 px-4 py-3 border-b">
              <RouterLink
                to="/"
                class="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
              >
                <ArrowLeft class="h-4 w-4" />
                <span>Back to Chat</span>
              </RouterLink>
            </div>
            <div class="flex items-center gap-3 px-4 py-3">
              <Settings class="h-5 w-5 text-muted-foreground" />
              <div>
                <h2 class="font-semibold text-sm">Settings</h2>
                <p class="text-xs text-muted-foreground">Manage your account</p>
              </div>
            </div>
          </SidebarHeader>

          <SidebarContent class="flex flex-col">
            <div class="p-2 flex-1 flex flex-col min-h-0">
              <div class="space-y-1">
                <RouterLink
                  v-for="item in settingsNavigation"
                  :key="item.name"
                  :to="item.href"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors hover:bg-accent/50"
                  active-class="bg-accent text-accent-foreground"
                >
                  <component :is="item.icon" class="h-4 w-4" />
                  <span>{{ item.name }}</span>
                </RouterLink>
              </div>
            </div>
          </SidebarContent>

          <SidebarFooter>
            <div class="flex items-center gap-2 px-2 py-2 bg-accent/50 rounded-lg mx-2 mb-2">
              <Avatar class="h-8 w-8">
                <AvatarImage :src="user.avatar" :alt="user.name" />
                <AvatarFallback>{{
                  user.name
                    .split(' ')
                    .map((n) => n[0])
                    .join('')
                }}</AvatarFallback>
              </Avatar>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate">{{ user.name }}</p>
                <p class="text-xs text-muted-foreground truncate">@{{ user.username }}</p>
              </div>
            </div>
          </SidebarFooter>
          <SidebarRail />
        </Sidebar>
      </Transition>

      <SidebarInset>
        <slot />
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
