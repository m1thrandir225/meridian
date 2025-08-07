<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar'
import { Settings, Plus, Waves } from 'lucide-vue-next'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RouterLink } from 'vue-router'
import { ref } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarFooter,
  SidebarRail,
} from '@/components/ui/sidebar'

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
  side: 'left',
})

const isNewChannelDialogOpen = ref(false)
const newChannelName = ref('')
const newChannelDescription = ref('')

// Sample server/workspace data
const server = {
  name: 'Meridian',
  avatar: '/server-avatar.png',
}

// Simplified channels data - only text channels
const channels = ref([
  { id: '1', name: 'general', unread: 0 },
  { id: '2', name: 'random', unread: 3 },
  { id: '3', name: 'development', unread: 0 },
  { id: '4', name: 'design', unread: 1 },
  { id: '5', name: 'announcements', unread: 0 },
  { id: '6', name: 'bug-reports', unread: 2 },
  { id: '7', name: 'feature-requests', unread: 0 },
  { id: '8', name: 'testing', unread: 1 },
  { id: '9', name: 'documentation', unread: 0 },
  { id: '10', name: 'off-topic', unread: 5 },
  { id: '11', name: 'help', unread: 0 },
  { id: '12', name: 'showcase', unread: 2 },
])

const user = {
  name: 'John Doe',
  username: 'johndoe',
  avatar: '/avatars/user.png',
}

const createChannel = () => {
  if (newChannelName.value.trim()) {
    const newChannel = {
      id: String(channels.value.length + 1),
      name: newChannelName.value.toLowerCase().replace(/\s+/g, '-'),
      unread: 0,
    }
    channels.value.push(newChannel)

    // Reset form
    newChannelName.value = ''
    newChannelDescription.value = ''
    isNewChannelDialogOpen.value = false
  }
}
</script>

<template>
  <Sidebar v-bind="props" class="border-r">
    <SidebarHeader>
      <div class="flex items-center gap-3 px-4 py-3 border-b">
        <Waves class="h-6 w-6" />
        <div class="flex-1 min-w-0">
          <h2 class="font-semibold text-sm truncate">{{ server.name }}</h2>
        </div>
        <Button variant="ghost" size="icon" class="h-6 w-6">
          <Settings class="h-4 w-4" />
        </Button>
      </div>
    </SidebarHeader>

    <SidebarContent class="flex flex-col">
      <div class="p-2 flex-1 flex flex-col min-h-0">
        <!-- Channels Section -->
        <div class="flex-1 flex flex-col min-h-0">
          <div class="flex items-center justify-between px-2 py-1 mb-2 flex-shrink-0">
            <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Channels
            </h3>
            <Dialog v-model:open="isNewChannelDialogOpen">
              <DialogTrigger as-child>
                <Button variant="ghost" size="icon" class="h-4 w-4">
                  <Plus class="h-3 w-3" />
                </Button>
              </DialogTrigger>
              <DialogContent class="sm:max-w-[425px]">
                <DialogHeader>
                  <DialogTitle>Create New Channel</DialogTitle>
                  <DialogDescription>
                    Add a new text channel to your server. Channel names should be lowercase and
                    descriptive.
                  </DialogDescription>
                </DialogHeader>
                <div class="grid gap-4 py-4">
                  <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="channel-name" class="text-right"> Name </Label>
                    <Input
                      id="channel-name"
                      v-model="newChannelName"
                      placeholder="e.g. general-chat"
                      class="col-span-3"
                    />
                  </div>
                  <div class="grid grid-cols-4 items-center gap-4">
                    <Label for="channel-description" class="text-right"> Description </Label>
                    <Input
                      id="channel-description"
                      v-model="newChannelDescription"
                      placeholder="Brief description (optional)"
                      class="col-span-3"
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button type="submit" @click="createChannel">Create Channel</Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
          <div class="space-y-0.5 overflow-y-auto flex-1 pr-2">
            <RouterLink
              v-for="channel in channels"
              :key="channel.id"
              :to="`/channel/${channel.id}`"
              class="flex items-center gap-2 px-2 py-1.5 rounded text-sm transition-colors hover:bg-accent/50"
              active-class="bg-accent text-accent-foreground"
            >
              <span class="flex-1 truncate">{{ channel.name }}</span>
              <Badge v-if="channel.unread > 0" variant="destructive" class="h-4 text-xs px-1">
                {{ channel.unread }}
              </Badge>
            </RouterLink>
          </div>
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
        <Button variant="ghost" size="icon" class="h-6 w-6">
          <Settings class="h-4 w-4" />
        </Button>
      </div>
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
