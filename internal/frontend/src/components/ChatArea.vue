<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getUserDisplayName, getUserInitials } from '@/lib/utils'
import websocketService from '@/services/websocket.service'
import { useAppearanceStore } from '@/stores/appearance'
import { useAuthStore } from '@/stores/auth'
import { useChannelStore } from '@/stores/channel'
import { useMessageStore } from '@/stores/message'
import type { Message } from '@/types/models/message'
import { Bot, Hash, Menu, Reply, Send, Smile, Users, X } from 'lucide-vue-next'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import ChannelInviteModal from './ChannelInviteModal.vue'
import ChannelSettings from './ChannelSettings.vue'
import MessageActionsPopup from './MessageActionsPopup.vue'
import MessageReactions from './MessageReactions.vue'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import TypingIndicator from './TypingIndicator.vue'

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

// State
const messagesContainer = ref<HTMLElement | null>(null)
const typingTimeout = ref<number | null>(null)
const isTyping = ref<boolean>(false)
const newMessage = ref<string>('')
const emojiPickerOpen = ref<boolean>(false)
const replyingTo = ref<Message | null>(null)
const messages = computed(() => messageStore.currentMessages)
const currentChannel = computed(() => channelStore.getCurrentChannel)
const isLoading = computed(() => messageStore.loading)
const inputRef = ref<{ inputRef: HTMLInputElement | null; focus: () => void } | null>(null)
const emojis = ref<string[]>([
  '👋',
  '👍',
  '👎',
  '🤔',
  '🤯',
  '🤮',
  '🤢',
  '🤪',
  '🤣',
  '😂',
  '😍',
  '😘',
  '😊',
  '😇',
  '😉',
  '😌',
  '😍',
  '😘',
  '😊',
  '😇',
  '😉',
  '😌',
  '🤑',
  '🤗',
  '🤓',
  '🤠',
  '🤡',
  '🤠',
])

const onEmojiSelect = (emoji: string) => {
  newMessage.value += emoji
  emojiPickerOpen.value = false
  setTimeout(() => {
    inputRef.value?.focus()
  }, 0)
}

// Computed styles based on appearance settings
const messageContainerClasses = () => {
  const baseClasses =
    'flex gap-3 group transition-all duration-200 rounded-lg hover:bg-accent/25 border border-transparent'
  const sizeClasses = appearanceStore.messageDisplayMode === 'compact' ? 'px-2 py-1' : 'px-3 py-2'

  return `${baseClasses} ${sizeClasses}`
}

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

    messageStore.sendMessage(newMessage.value, currentChannel.value.id, replyingTo.value?.id)

    // Clear input
    newMessage.value = ''
    replyingTo.value = null

    stopTyping()
  }
}

const startReply = (message: Message) => {
  replyingTo.value = message

  nextTick(() => {
    const input = document.querySelector(
      'input[placeholder="Type a message..."]',
    ) as HTMLInputElement
    if (input) {
      input.focus()
    }
  })
}

const cancelReply = () => {
  replyingTo.value = null
}

const getParentMessage = (message: Message): Message | null => {
  if (!message.parent_message_id) return null

  return messages.value.find((m) => m.id === message.parent_message_id) ?? null
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

const handleReply = (message: Message) => {
  startReply(message)
}
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
        <div v-if="currentChannel" class="px-2">
          <ChannelInviteModal :channel-id="currentChannel.id" :channel-name="currentChannel.name" />
        </div>
        <ChannelSettings v-if="currentChannel" :channel="currentChannel" />
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

      <div v-else :class="messageSpacing">
        <MessageActionsPopup
          v-for="message in messages"
          :key="message.id"
          :message="message"
          @reply="handleReply"
        >
          <div :class="messageContainerClasses()" class="flex gap-3 relative w-full">
            <!-- Avatar -->
            <Avatar :class="avatarSize" v-if="message.sender_user_id">
              <AvatarFallback v-if="message.sender_user_id">
                {{
                  message.sender_user_id === authStore.user?.id
                    ? `${getUserInitials(authStore.userDisplayName())}`
                    : getUserInitials(
                        `${message.sender_user?.first_name} ${message.sender_user?.last_name}`,
                      )
                }}
              </AvatarFallback>
            </Avatar>
            <Avatar :class="avatarSize" v-else-if="message.integration_id">
              <AvatarFallback>
                <Bot />
              </AvatarFallback>
            </Avatar>

            <!-- Message Content -->
            <div class="flex-1 min-w-0">
              <div v-if="message.parent_message_id" class="mb-1">
                <div class="flex items-center gap-1 text-xs text-muted-foreground">
                  <Reply class="h-3 w-3" />
                  <span>Replying to</span>
                  <span class="font-medium" v-if="getParentMessage(message)?.sender_user_id">
                    {{
                      getParentMessage(message)?.sender_user_id === authStore.user?.id
                        ? 'You'
                        : getUserDisplayName(
                            getParentMessage(message)?.sender_user?.first_name ?? '',
                            getParentMessage(message)?.sender_user?.last_name ?? '',
                          )
                    }}
                  </span>
                  <span class="font-medium" v-else-if="getParentMessage(message)?.integration_id">
                    {{ getParentMessage(message)?.integration_bot?.service_name }}
                  </span>
                </div>
                <div class="text-xs text-muted-foreground truncate max-w-xs">
                  {{ getParentMessage(message)?.content_text }}
                </div>
              </div>

              <div class="flex items-center gap-2 mb-1">
                <span
                  class="font-semibold text-primary text-[10px]"
                  v-if="message.sender_user_id && message.sender_user_id !== authStore.user?.id"
                >
                  @{{ message.sender_user?.username }}
                </span>
                <span class="font-medium text-sm" v-if="message.sender_user_id">
                  {{
                    message.sender_user_id === authStore.user?.id
                      ? 'You'
                      : getUserDisplayName(
                          message.sender_user?.first_name ?? '',
                          message.sender_user?.last_name ?? '',
                        )
                  }}
                </span>
                <span class="font-medium text-sm" v-else-if="message.integration_id">
                  {{ message.integration_bot?.service_name }}
                </span>
                <span class="text-xs text-muted-foreground">
                  {{ formatTime(message.created_at) }}
                </span>
              </div>

              <div :class="messageTextClasses">
                {{ message.content_text }}
              </div>

              <MessageReactions :message="message" />
            </div>
          </div>
        </MessageActionsPopup>
      </div>
    </div>

    <!-- Typing Indicator -->
    <TypingIndicator v-if="currentChannel?.id" :channel-id="currentChannel.id" />

    <!-- Reply Area -->
    <div v-if="replyingTo" class="px-4 py-2 bg-muted/50 border-t border-border">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2 text-sm">
          <Reply class="h-4 w-4 text-muted-foreground" />
          <span class="text-muted-foreground">Replying to</span>
          <span class="font-medium">
            {{
              replyingTo.sender_user_id === authStore.user?.id
                ? 'You'
                : getUserDisplayName(
                    replyingTo.sender_user?.first_name ?? '',
                    replyingTo.sender_user?.last_name ?? '',
                  )
            }}
          </span>
        </div>
        <Button variant="ghost" size="sm" class="h-6 w-6 p-0" @click="cancelReply">
          <X class="h-3 w-3" />
        </Button>
      </div>
      <div class="text-sm text-muted-foreground truncate max-w-md mt-1">
        {{ replyingTo.content_text }}
      </div>
    </div>

    <!-- Input Area -->
    <div class="p-4 border-t">
      <div class="flex gap-2">
        <div class="flex-1 relative">
          <Input
            ref="inputRef"
            v-model="newMessage"
            placeholder="Type a message..."
            class="pr-20"
            @input="handleInputChange"
            @keypress="handleKeyPress"
          />

          <div class="absolute right-2 top-1/2 transform -translate-y-1/2 flex gap-1">
            <Popover v-model:open="emojiPickerOpen">
              <PopoverTrigger as-child>
                <Button variant="ghost" size="icon" class="h-6 w-6">
                  <Smile class="h-3 w-3" />
                </Button>
              </PopoverTrigger>
              <PopoverContent side="top" align="center" class="p-2 w-64">
                <div class="grid grid-cols-8 gap-1">
                  <button
                    v-for="emoji in emojis"
                    :key="emoji"
                    class="h-7 w-7 flex items-center justify-center rounded-md text-2xl hover:bg-accent"
                    @click="onEmojiSelect(emoji)"
                  >
                    {{ emoji }}
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          </div>
        </div>

        <Button @click="sendMessage" :disabled="!newMessage.trim()" size="icon" class="h-9 w-9">
          <Send class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>
