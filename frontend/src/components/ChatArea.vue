<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted } from 'vue'
import { Hash, Users, Send, Smile, Menu } from 'lucide-vue-next'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAppearanceStore } from '@/stores/appearance'

// Appearance store
const appearanceStore = useAppearanceStore()

// Chat scroll management
const messagesContainer = ref<HTMLElement | null>(null)
const isUserInteracting = ref(false)
const scrollToBottomTimeout = ref<number | null>(null)

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

// Computed styles based on appearance settings
const messageContainerClasses = computed(() => {
  const baseClasses = 'flex gap-3 group transition-all duration-200 rounded-lg hover:bg-accent/25'
  const sizeClasses = appearanceStore.messageDisplayMode === 'compact' ? 'px-2 py-1' : 'px-3 py-2'

  return `${baseClasses} ${sizeClasses}`
})

const messageTextClasses = computed(() => {
  const baseClasses = 'leading-relaxed mb-2'
  const sizeClasses = appearanceStore.messageDisplayMode === 'compact' ? 'text-sm' : 'text-base'

  return `${baseClasses} ${sizeClasses}`
})

const avatarSize = computed(() => {
  return appearanceStore.messageDisplayMode === 'compact' ? 'h-6 w-6' : 'h-8 w-8'
})

const messageSpacing = computed(() => {
  return appearanceStore.messageDisplayMode === 'compact' ? 'space-y-1' : 'space-y-4'
})

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

    // Force scroll to bottom when user sends a message
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
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

// Track hover state for messages
const hoveredMessageId = ref<string | null>(null)

const setHoveredMessage = (messageId: string | null) => {
  hoveredMessageId.value = messageId
}

const getMessageStyle = (messageId: string) => {
  const isHovered = hoveredMessageId.value === messageId
  if (!isHovered)
    return {
      border: '2px solid transparent',
    }

  return {
    border: `2px solid hsl(${appearanceStore.accentColorClass})`,
    marginLeft: '4px',
    paddingLeft: '8px',
  }
}

// Auto-scroll functionality
const scrollToBottom = () => {
  if (messagesContainer.value && !isUserInteracting.value) {
    nextTick(() => {
      messagesContainer.value!.scrollTop = messagesContainer.value!.scrollHeight
    })
  }
}

const handleScrollInteraction = () => {
  isUserInteracting.value = true

  if (scrollToBottomTimeout.value) {
    clearTimeout(scrollToBottomTimeout.value)
  }

  // Resume auto-scroll after user stops scrolling for 3 seconds
  scrollToBottomTimeout.value = window.setTimeout(() => {
    isUserInteracting.value = false
  }, 3000)
}

const handleMessageHover = (hovering: boolean) => {
  if (hovering) {
    isUserInteracting.value = true
  } else {
    // Short delay before resuming auto-scroll when user stops hovering
    if (scrollToBottomTimeout.value) {
      clearTimeout(scrollToBottomTimeout.value)
    }
    scrollToBottomTimeout.value = window.setTimeout(() => {
      isUserInteracting.value = false
    }, 1000)
  }
}

// Watch for new messages and auto-scroll
watch(
  () => messages.value.length,
  () => {
    nextTick(() => {
      scrollToBottom()
    })
  },
)

// Initial scroll to bottom on mount
onMounted(() => {
  nextTick(() => {
    scrollToBottom()
  })
})
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
    <div
      ref="messagesContainer"
      class="flex-1 overflow-y-auto p-4"
      :class="messageSpacing"
      @scroll="handleScrollInteraction"
      @wheel="handleScrollInteraction"
    >
      <div
        v-for="message in messages"
        :key="message.id"
        :class="messageContainerClasses"
        :style="getMessageStyle(message.id)"
        @mouseenter="(setHoveredMessage(message.id), handleMessageHover(true))"
        @mouseleave="(setHoveredMessage(null), handleMessageHover(false))"
      >
        <Avatar :class="`${avatarSize} mt-1`">
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
            <span
              :class="
                appearanceStore.messageDisplayMode === 'compact'
                  ? 'font-medium text-xs'
                  : 'font-medium text-sm'
              "
              >{{ message.author.name }}</span
            >
            <span
              :class="
                appearanceStore.messageDisplayMode === 'compact'
                  ? 'text-xs text-muted-foreground'
                  : 'text-xs text-muted-foreground'
              "
              >{{ formatTime(message.timestamp) }}</span
            >
          </div>
          <div :class="messageTextClasses">
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
          <div class="opacity-0 group-hover:opacity-100 transition-opacity w-auto">
            <div class="flex gap-1 mt-1 px-4 py-2 w-fit border border-primary/50 rounded-xl">
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
