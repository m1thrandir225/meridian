<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { Hash, Users, Send, Smile, Menu, Paperclip } from 'lucide-vue-next'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAppearanceStore } from '@/stores/appearance'
import { useChannelStore } from '@/stores/channel'
import { useMessageStore } from '@/stores/message'
import { useAuthStore } from '@/stores/auth'
import websocketService from '@/services/websocket.service'
import { getUserInitials, getUserDisplayName } from '@/lib/utils'

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

// Appearance store
const appearanceStore = useAppearanceStore()
const messageStore = useMessageStore()
const channelStore = useChannelStore()
const authStore = useAuthStore()

// Chat scroll management
const messagesContainer = ref<HTMLElement | null>(null)
const isUserInteracting = ref(false)
const scrollToBottomTimeout = ref<number | null>(null)
const typingTimeout = ref<number | null>(null)
const isTyping = ref<boolean>(false)
const newMessage = ref<string>('')

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

const messages = computed(() => messageStore.currentMessages)
const currentChannel = computed(() => channelStore.getCurrentChannel)
const isLoading = computed(() => messageStore.loading)

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

watch(
  messages,
  () => {
    scrollToBottom()
  },
  { deep: true },
)

const sendMessage = () => {
  if (newMessage.value.trim() && currentChannel.value?.id) {
    console.log('Sending message:', newMessage.value)

    messageStore.sendMessage(newMessage.value, currentChannel.value.id)

    // Clear input
    newMessage.value = ''

    stopTyping()
  }
}

const startTyping = () => {
  if (!currentChannel.value?.id || isTyping.value) return

  isTyping.value = true
  messageStore.startTyping(currentChannel.value.id)
}

const stopTyping = () => {
  if (!currentChannel.value?.id) return

  isTyping.value = false
  messageStore.stopTyping(currentChannel.value.id)
}

const handleInputChange = () => {
  startTyping()

  // Clear existing timeout
  if (typingTimeout.value) {
    clearTimeout(typingTimeout.value)
  }

  // Set new timeout to stop typing after 2 seconds
  typingTimeout.value = setTimeout(() => {
    stopTyping()
  }, 2000)
}

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

onUnmounted(() => {
  if (typingTimeout.value) {
    clearTimeout(typingTimeout.value)
  }
  stopTyping()
})

onMounted(() => {
  if (!websocketService.isConnected()) {
    websocketService.connect()
  }
})

const formatTime = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

// const toggleReaction = (messageId: string, emoji: string) => {
//   const message = messages.value.find((m) => m.id === messageId)
//   if (!message) return

//   const existingReaction = message.reactions.find((r) => r.emoji === emoji)

//   if (existingReaction) {
//     if (existingReaction.userReacted) {
//       existingReaction.count--
//       existingReaction.userReacted = false
//       if (existingReaction.count === 0) {
//         message.reactions = message.reactions.filter((r) => r.emoji !== emoji)
//       }
//     } else {
//       existingReaction.count++
//       existingReaction.userReacted = true
//     }
//   } else {
//     message.reactions.push({
//       emoji,
//       count: 1,
//       userReacted: true,
//     })
//   }
// }

//const quickReactions = ['👍', '❤️', '😂', '😮', '😢', '😡']

// Track hover state for messages
//const hoveredMessageId = ref<string | null>(null)

// const setHoveredMessage = (messageId: string | null) => {
//   hoveredMessageId.value = messageId
// }

// const getMessageStyle = (messageId: string) => {
//   const isHovered = hoveredMessageId.value === messageId
//   if (!isHovered)
//     return {
//       border: '2px solid transparent',
//     }

//   return {
//     border: `2px solid hsl(${appearanceStore.accentColorClass})`,
//     marginLeft: '4px',
//     paddingLeft: '8px',
//   }
// }

// Auto-scroll functionality
// const scrollToBottom = () => {
//   if (messagesContainer.value && !isUserInteracting.value) {
//     nextTick(() => {
//       messagesContainer.value!.scrollTop = messagesContainer.value!.scrollHeight
//     })
//   }
// }

// const handleScrollInteraction = () => {
//   isUserInteracting.value = true

//   if (scrollToBottomTimeout.value) {
//     clearTimeout(scrollToBottomTimeout.value)
//   }

//   // Resume auto-scroll after user stops scrolling for 3 seconds
//   scrollToBottomTimeout.value = window.setTimeout(() => {
//     isUserInteracting.value = false
//   }, 3000)
// }

// const handleMessageHover = (hovering: boolean) => {
//   if (hovering) {
//     isUserInteracting.value = true
//   } else {
//     // Short delay before resuming auto-scroll when user stops hovering
//     if (scrollToBottomTimeout.value) {
//       clearTimeout(scrollToBottomTimeout.value)
//     }
//     scrollToBottomTimeout.value = window.setTimeout(() => {
//       isUserInteracting.value = false
//     }, 1000)
//   }
// }

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
          <h1 class="font-semibold">{{ channelStore.currentChannel?.name }}</h1>
          <p class="text-sm text-muted-foreground">{{ channelStore.currentChannel?.topic }}</p>
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
    <div class="flex-1 overflow-y-auto p-4" ref="messagesContainer">
      <div v-if="isLoading" class="flex justify-center items-center h-full">
        <div class="text-muted-foreground">Loading messages...</div>
      </div>

      <div v-else-if="messages.length === 0" class="flex justify-center items-center h-full">
        <div class="text-center text-muted-foreground">
          <p>No messages yet</p>
          <p class="text-sm">Start the conversation!</p>
        </div>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="message in messages"
          :key="message.id"
          :class="messageContainerClasses"
          class="flex gap-3"
        >
          <!-- Avatar -->
          <Avatar class="h-8 w-8">
            <AvatarFallback>
              {{
                message.sender_user_id === authStore.user?.user_id
                  ? `${getUserInitials(authStore.userDisplayName())}`
                  : getUserInitials(
                      `${message.sender_user?.first_name} ${message.sender_user?.last_name}`,
                    )
              }}
            </AvatarFallback>
          </Avatar>

          <!-- Message Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span
                class="font-semibold text-primary text-[10px]"
                v-if="message.sender_user_id !== authStore.user?.user_id"
              >
                @{{ message.sender_user?.username }}
              </span>
              <span class="font-medium text-sm">
                {{
                  message.sender_user_id === authStore.user?.user_id
                    ? 'You'
                    : getUserDisplayName(
                        message.sender_user?.first_name ?? '',
                        message.sender_user?.last_name ?? '',
                      )
                }}
              </span>
              <span class="text-xs text-muted-foreground">
                {{ formatTime(message.created_at) }}
              </span>
            </div>

            <div class="text-sm">
              {{ message.content_text }}
            </div>

            <!-- Reactions (if any) -->
            <!-- <div v-if="message.reactions && message.reactions.length > 0" class="flex gap-1 mt-2">
              <Badge
                v-for="reaction in message.reactions"
                :key="reaction.type"
                variant="secondary"
                class="text-xs"
              >
                {{ reaction.type }} {{ reaction.count }}
              </Badge>
            </div> -->
          </div>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="p-4 border-t">
      <div class="flex gap-2">
        <Button variant="ghost" size="icon" class="h-9 w-9">
          <Paperclip class="h-4 w-4" />
        </Button>

        <div class="flex-1 relative">
          <Input
            v-model="newMessage"
            placeholder="Type a message..."
            class="pr-20"
            @input="handleInputChange"
            @keypress="handleKeyPress"
          />

          <div class="absolute right-2 top-1/2 transform -translate-y-1/2 flex gap-1">
            <Button variant="ghost" size="icon" class="h-6 w-6">
              <Smile class="h-3 w-3" />
            </Button>
          </div>
        </div>

        <Button @click="sendMessage" :disabled="!newMessage.trim()" size="icon" class="h-9 w-9">
          <Send class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>
