<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import type { SidebarProps } from '@/components/ui/sidebar'
import { Bot, Plus, Users } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'

import { Sidebar, SidebarContent } from '@/components/ui/sidebar'
import { getUserDisplayName, getUserInitials } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth'
import { useChannelStore } from '@/stores/channel'
import { computed } from 'vue'

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
  side: 'right',
})

const channelStore = useChannelStore()
const authStore = useAuthStore()

const currentUser = computed(() => authStore.user)
const currentChannel = computed(() => channelStore.currentChannel)

const bots = computed(() => currentChannel.value?.bots)
const activeBots = computed(() => bots.value?.filter((bot) => !bot.is_revoked))
const members = computed(() => currentChannel.value?.members)
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
          <div v-if="currentUser?.id === currentChannel?.creator_user_id" class="px-2">
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
                Active Bots
              </h3>
              <Badge variant="secondary" class="text-xs">
                {{ activeBots?.length }}
              </Badge>
            </div>
            <div v-if="activeBots && activeBots?.length > 0" class="space-y-1">
              <div
                v-for="bot in activeBots"
                :key="bot.id"
                class="flex items-center gap-3 rounded-lg p-2 hover:bg-accent transition-colors cursor-pointer"
              >
                <div class="relative">
                  <Avatar class="h-8 w-8">
                    <AvatarFallback>{{ getUserInitials(bot.service_name) }}</AvatarFallback>
                  </Avatar>
                  <div
                    class="absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full border-2 border-background bg-blue-500 flex items-center justify-center"
                  >
                    <Bot class="h-2 w-2 text-white" />
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium truncate">{{ bot.service_name }}</p>
                  <p class="text-xs text-muted-foreground truncate">
                    {{ new Date(bot.created_at).toLocaleDateString() }}
                  </p>
                </div>
              </div>
            </div>
            <p v-else class="text-sm text-muted-foreground">No active bots.</p>
          </div>

          <!-- Members Section (Second) -->
          <div>
            <div class="flex items-center justify-between px-2 py-1 mb-2">
              <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                Members
              </h3>
              <Badge variant="secondary" class="text-xs">
                {{ members?.length }}
              </Badge>
            </div>
            <div class="space-y-1">
              <div
                v-for="member in members"
                :key="member.id"
                class="flex items-center gap-3 rounded-lg p-2 hover:bg-accent transition-colors cursor-pointer"
              >
                <Avatar class="h-8 w-8">
                  <AvatarFallback>{{
                    getUserInitials(getUserDisplayName(member.first_name, member.last_name))
                  }}</AvatarFallback>
                </Avatar>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1">
                    <p class="text-[13px] font-medium truncate">@{{ member.username }}</p>
                  </div>
                  <p class="text-xs text-muted-foreground truncate">
                    {{ getUserDisplayName(member.first_name, member.last_name) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </SidebarContent>
  </Sidebar>
</template>
