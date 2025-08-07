<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Users, Crown, Bot, Plus } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'

import { Sidebar, SidebarContent, SidebarRail } from '@/components/ui/sidebar'

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
  side: 'right',
})

// Current user data (for checking if owner)
const currentUser = {
  id: '1',
  role: 'owner', // This would come from your auth system
}

// Sample members data - simplified without status
const members = [
  {
    id: '1',
    name: 'John Doe',
    username: 'johndoe',
    avatar: '/avatars/01.png',
    role: 'owner',
  },
  {
    id: '2',
    name: 'Jane Smith',
    username: 'janesmith',
    avatar: '/avatars/02.png',
    role: 'member',
  },
  {
    id: '3',
    name: 'Bob Johnson',
    username: 'bobjohnson',
    avatar: '/avatars/03.png',
    role: 'member',
  },
  {
    id: '4',
    name: 'Alice Brown',
    username: 'alicebrown',
    avatar: '/avatars/04.png',
    role: 'member',
  },
  {
    id: '5',
    name: 'Charlie Wilson',
    username: 'charlie',
    avatar: '/avatars/05.png',
    role: 'member',
  },
  {
    id: '6',
    name: 'Diana Prince',
    username: 'diana',
    avatar: '/avatars/06.png',
    role: 'member',
  },
  {
    id: '7',
    name: 'Eddie Murphy',
    username: 'eddie',
    avatar: '/avatars/07.png',
    role: 'member',
  },
  {
    id: '8',
    name: 'Fiona Green',
    username: 'fiona',
    avatar: '/avatars/08.png',
    role: 'member',
  },
]

// Sample bots data
const bots = [
  {
    id: 'bot1',
    name: 'ModerationBot',
    username: 'modbot',
    avatar: '/avatars/bot1.png',
  },
  {
    id: 'bot2',
    name: 'WelcomeBot',
    username: 'welcome',
    avatar: '/avatars/bot2.png',
  },
  {
    id: 'bot3',
    name: 'MusicBot',
    username: 'music',
    avatar: '/avatars/bot3.png',
  },
  {
    id: 'bot4',
    name: 'AnnouncementBot',
    username: 'announce',
    avatar: '/avatars/bot4.png',
  },
  {
    id: 'bot5',
    name: 'SecurityBot',
    username: 'security',
    avatar: '/avatars/bot5.png',
  },
  {
    id: 'bot6',
    name: 'AnalyticsBot',
    username: 'analytics',
    avatar: '/avatars/bot6.png',
  },
]

const getRoleIcon = (role: string) => {
  if (role === 'owner') return Crown
  return null
}
</script>

<template>
  <Sidebar v-bind="props" class="border-l">
    <SidebarContent class="flex flex-col">
      <div class="p-2 flex-1 flex flex-col min-h-0">
        <div class="flex items-center gap-2 px-2 py-1 mb-3 flex-shrink-0">
          <Users class="h-4 w-4" />
          <span class="font-semibold text-sm">Channel Info</span>
        </div>

        <!-- Combined scrollable area for bots and members -->
        <div class="space-y-4 overflow-y-auto flex-1 pr-2">
          <!-- Bot Registration Button (Only for owners) -->
          <div v-if="currentUser.role === 'owner'" class="px-2">
            <RouterLink
              to="/bot-registration"
              class="flex items-center gap-2 w-full px-3 py-2 text-sm bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors"
            >
              <Plus class="h-4 w-4" />
              <span>Register New Bot</span>
            </RouterLink>
          </div>

          <!-- Bots Section (First) -->
          <div>
            <div class="flex items-center justify-between px-2 py-1 mb-2">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                Bots
              </h3>
              <Badge variant="secondary" class="text-xs">
                {{ bots.length }}
              </Badge>
            </div>
            <div class="space-y-1">
              <div
                v-for="bot in bots"
                :key="bot.id"
                class="flex items-center gap-3 rounded-lg p-2 hover:bg-accent transition-colors cursor-pointer"
              >
                <div class="relative">
                  <Avatar class="h-8 w-8">
                    <AvatarImage :src="bot.avatar" :alt="bot.name" />
                    <AvatarFallback>{{
                      bot.name
                        .split(' ')
                        .map((n) => n[0])
                        .join('')
                    }}</AvatarFallback>
                  </Avatar>
                  <div
                    class="absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full border-2 border-background bg-blue-500 flex items-center justify-center"
                  >
                    <Bot class="h-2 w-2 text-white" />
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium truncate">{{ bot.name }}</p>
                  <p class="text-xs text-muted-foreground truncate">@{{ bot.username }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Members Section (Second) -->
          <div>
            <div class="flex items-center justify-between px-2 py-1 mb-2">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                Members
              </h3>
              <Badge variant="secondary" class="text-xs">
                {{ members.length }}
              </Badge>
            </div>
            <div class="space-y-1">
              <div
                v-for="member in members"
                :key="member.id"
                class="flex items-center gap-3 rounded-lg p-2 hover:bg-accent transition-colors cursor-pointer"
              >
                <Avatar class="h-8 w-8">
                  <AvatarImage :src="member.avatar" :alt="member.name" />
                  <AvatarFallback>{{
                    member.name
                      .split(' ')
                      .map((n) => n[0])
                      .join('')
                  }}</AvatarFallback>
                </Avatar>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1">
                    <p class="text-sm font-medium truncate">{{ member.name }}</p>
                    <component
                      v-if="getRoleIcon(member.role)"
                      :is="getRoleIcon(member.role)"
                      class="h-3 w-3 text-amber-500"
                    />
                  </div>
                  <p class="text-xs text-muted-foreground truncate">@{{ member.username }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </SidebarContent>
    <SidebarRail />
  </Sidebar>
</template>
