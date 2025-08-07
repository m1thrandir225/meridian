<script setup lang="ts">
import { ref } from 'vue'
import { Hash, Users, Send, Smile, Menu } from 'lucide-vue-next'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

// Props and emits
interface Props {
  isChannelsSidebarOpen?: boolean
  isMembersSidebarOpen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isChannelsSidebarOpen: true,
  isMembersSidebarOpen: true,
})

const emit = defineEmits<{
  toggleChannelsSidebar: []
  toggleMembersSidebar: []
}>()

const currentChannel = {
  id: '1',
  name: 'general',
  description: 'General discussion channel',
  type: 'text',
}

const messages = ref([
  {
    id: '1',
    author: {
      name: 'John Doe',
      username: 'johndoe',
      avatar: '/avatars/01.png',
    },
    content: 'Hey everyone! Welcome to the general channel 👋',
    timestamp: '2024-01-15T10:30:00Z',
    isOwn: false,
    reactions: [
      { emoji: '👋', count: 3, userReacted: false },
      { emoji: '🎉', count: 1, userReacted: true },
    ],
  },
  {
    id: '2',
    author: {
      name: 'Jane Smith',
      username: 'janesmith',
      avatar: '/avatars/02.png',
    },
    content: 'Thanks for setting this up! This looks great.',
    timestamp: '2024-01-15T10:32:00Z',
    isOwn: false,
    reactions: [{ emoji: '👍', count: 2, userReacted: false }],
  },
  {
    id: '3',
    author: {
      name: 'You',
      username: 'you',
      avatar: '/avatars/user.png',
    },
    content: 'Glad you like it! Feel free to share any feedback.',
    timestamp: '2024-01-15T10:35:00Z',
    isOwn: true,
    reactions: [],
  },
])

const newMessage = ref('')

const formatTime = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

const sendMessage = () => {
  if (newMessage.value.trim()) {
    console.log('Sending message:', newMessage.value)
    messages.value.push({
      id: String(messages.value.length + 1),
      author: {
        name: 'You',
        username: 'you',
        avatar: '/avatars/user.png',
      },
      content: newMessage.value,
      timestamp: new Date().toISOString(),
      isOwn: true,
      reactions: [],
    })
    newMessage.value = ''
  }
}

const toggleReaction = (messageId: string, emoji: string) => {
  const message = messages.value.find((m) => m.id === messageId)
  if (!message) return

  const existingReaction = message.reactions.find((r) => r.emoji === emoji)

  if (existingReaction) {
    if (existingReaction.userReacted) {
      existingReaction.count--
      existingReaction.userReacted = false
      if (existingReaction.count === 0) {
        message.reactions = message.reactions.filter((r) => r.emoji !== emoji)
      }
    } else {
      existingReaction.count++
      existingReaction.userReacted = true
    }
  } else {
    message.reactions.push({
      emoji,
      count: 1,
      userReacted: true,
    })
  }
}

const quickReactions = ['👍', '❤️', '😂', '😮', '😢', '😡']
</script>

<template>
  <div class="flex flex-col h-full w-full">
    <!-- Channel Header -->
    <header
      class="flex items-center justify-between px-4 py-3 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
    >
      <div class="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon"
          @click="emit('toggleChannelsSidebar')"
          :class="props.isChannelsSidebarOpen ? 'bg-accent' : ''"
        >
          <Menu class="h-4 w-4" />
        </Button>
        <Hash class="h-5 w-5 text-muted-foreground" />
        <div>
          <h1 class="font-semibold">{{ currentChannel.name }}</h1>
          <p class="text-sm text-muted-foreground">{{ currentChannel.description }}</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          @click="emit('toggleMembersSidebar')"
          :class="props.isMembersSidebarOpen ? 'bg-accent' : ''"
        >
          <Users class="h-4 w-4" />
        </Button>
      </div>
    </header>

    <!-- Messages Area -->
    <div class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-for="message in messages" :key="message.id" class="flex gap-3 group">
        <Avatar class="h-8 w-8 mt-1">
          <AvatarImage :src="message.author.avatar" :alt="message.author.name" />
          <AvatarFallback>{{
            message.author.name
              .split(' ')
              .map((n) => n[0])
              .join('')
          }}</AvatarFallback>
        </Avatar>
        <div class="flex-1 min-w-0">
          <div class="flex items-baseline gap-2 mb-1">
            <span class="font-medium text-sm">{{ message.author.name }}</span>
            <span class="text-xs text-muted-foreground">{{ formatTime(message.timestamp) }}</span>
          </div>
          <div class="text-sm leading-relaxed mb-2">
            {{ message.content }}
          </div>

          <!-- Reactions Display -->
          <div v-if="message.reactions.length > 0" class="flex flex-wrap gap-1 mb-2">
            <button
              v-for="reaction in message.reactions"
              :key="reaction.emoji"
              @click="toggleReaction(message.id, reaction.emoji)"
              :class="[
                'inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs border transition-colors',
                reaction.userReacted
                  ? 'bg-blue-100 border-blue-300 text-blue-700'
                  : 'bg-muted border-border hover:bg-accent',
              ]"
            >
              <span>{{ reaction.emoji }}</span>
              <span>{{ reaction.count }}</span>
            </button>
          </div>

          <!-- Quick Reaction Buttons (appear on hover) -->
          <div class="opacity-0 group-hover:opacity-100 transition-opacity">
            <div class="flex gap-1 mt-1">
              <button
                v-for="emoji in quickReactions"
                :key="emoji"
                @click="toggleReaction(message.id, emoji)"
                class="w-6 h-6 rounded text-xs hover:bg-accent transition-colors flex items-center justify-center"
                :title="`React with ${emoji}`"
              >
                {{ emoji }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Message Input -->
    <div class="p-4 border-t">
      <div class="relative flex items-center gap-2 bg-muted/50 rounded-lg p-3">
        <Input
          v-model="newMessage"
          placeholder="Type a message..."
          class="flex-1 border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
          @keydown.enter="sendMessage"
        />
        <Button variant="ghost" size="icon" class="h-6 w-6">
          <Smile class="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          class="h-6 w-6"
          @click="sendMessage"
          :disabled="!newMessage.trim()"
        >
          <Send class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>
