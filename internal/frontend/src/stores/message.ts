import channelService from '@/services/channel.service'
import websocketService from '@/services/websocket.service'
import { useAuthStore } from '@/stores/auth'
import type { Message } from '@/types/models/message'
import type { Reaction } from '@/types/models/reaction'
import type {
  IncomingMessagePayload,
  IncomingReactionPayload,
  TypingPayload,
} from '@/types/websocket'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const loading = ref(false)
  const currentChannelId = ref<string | null>(null)
  const typingUsers = ref<Map<string, Set<string>>>(new Map()) // channelId -> Set of usernames
  const authStore = useAuthStore()

  const currentMessages = computed(() => {
    if (!currentChannelId.value) return []
    return messages.value.filter((message) => message.channel_id === currentChannelId.value)
  })

  async function fetchMessages(channelId: string) {
    loading.value = true
    try {
      const fetchedMessages = await channelService.getMessages(channelId)
      const existingIds = new Set(messages.value.map((m) => m.id))
      const newMessages = fetchedMessages.filter((m) => !existingIds.has(m.id))
      messages.value = [...messages.value, ...newMessages]
      currentChannelId.value = channelId
    } catch (error) {
      console.error('Error fetching messages:', error)
    } finally {
      loading.value = false
    }
  }

  function addMessage(message: Message) {
    const existingIndex = messages.value.findIndex((m) => m.id === message.id)
    if (existingIndex === -1) {
      messages.value.push(message)
    } else {
      // Update existing message
      messages.value[existingIndex] = message
    }
  }

  function addMessageFromWebSocket(payload: IncomingMessagePayload) {
    console.log('payload', payload)
    const message: Message = {
      id: payload.id,
      channel_id: payload.channel_id,
      sender_user_id: payload.sender_user_id,
      integration_id: payload.integration_id,
      content_text: payload.content,
      parent_message_id: payload.parent_message_id,
      created_at: payload.timestamp,
      sender_user: payload.sender_user,
      integration_bot: payload.integration_bot,
    }

    addMessage(message)
  }

  function addReactionToMessage(messageId: string, reactionPayload: IncomingReactionPayload) {
    console.log('Adding reaction to message:', messageId, reactionPayload)
    const messageIndex = messages.value.findIndex((m) => m.id === messageId)
    if (messageIndex === -1) {
      console.log('Message not found for reaction:', messageId)
      return
    }

    const message = messages.value[messageIndex]
    if (!message.reactions) {
      message.reactions = []
    }

    // Convert payload to Reaction type
    const reaction: Reaction = {
      id: reactionPayload.id,
      message_id: reactionPayload.message_id,
      user_id: reactionPayload.user_id,
      reaction_type: reactionPayload.reaction_type,
      timestamp: reactionPayload.timestamp,
    }

    // Check if user already has this reaction
    const existingReactionIndex = message.reactions.findIndex(
      (r) => r.user_id === reaction.user_id && r.reaction_type === reaction.reaction_type,
    )

    if (existingReactionIndex === -1) {
      message.reactions.push(reaction)
      console.log('Reaction added, new count:', message.reactions.length)
    } else {
      console.log('Reaction already exists, skipping')
    }
  }

  function removeReactionFromMessage(messageId: string, userId: string, reactionType: string) {
    console.log('Removing reaction from message:', messageId, userId, reactionType)
    const messageIndex = messages.value.findIndex((m) => m.id === messageId)
    if (messageIndex === -1) {
      console.log('Message not found for reaction removal:', messageId)
      return
    }

    const message = messages.value[messageIndex]
    if (!message.reactions) {
      console.log('No reactions to remove')
      return
    }

    const beforeCount = message.reactions.length
    message.reactions = message.reactions.filter(
      (r) => !(r.user_id === userId && r.reaction_type === reactionType),
    )
    const afterCount = message.reactions.length
    console.log(`Reaction removed: ${beforeCount} -> ${afterCount}`)
  }

  function addReaction(messageId: string, reactionType: string) {
    if (!currentChannelId.value || !authStore.user) return

    console.log('Sending add reaction:', messageId, reactionType)
    websocketService.sendAddReaction(currentChannelId.value, messageId, reactionType)
  }

  function removeReaction(messageId: string, reactionType: string) {
    if (!currentChannelId.value || !authStore.user) return

    console.log('Sending remove reaction:', messageId, reactionType)
    websocketService.sendRemoveReaction(currentChannelId.value, messageId, reactionType)
  }

  function clearMessages() {
    messages.value = []
    currentChannelId.value = null
  }

  function sendMessage(content: string, channelId: string, parentMessageId?: string) {
    websocketService.sendMessage({
      content,
      channel_id: channelId,
      parent_message_id: parentMessageId,
    })
  }

  function startTyping(channelId: string) {
    websocketService.sendTypingStart(channelId)
  }

  function stopTyping(channelId: string) {
    websocketService.sendTypingStop(channelId)
  }

  function handleTypingStart(payload: TypingPayload) {
    if (!payload.channel_id || !payload.user_id || payload.user_id === authStore.user?.id) {
      return // Don't show own typing indicator
    }

    const username = payload.username || payload.user?.username || 'Unknown User'

    if (!typingUsers.value.has(payload.channel_id)) {
      typingUsers.value.set(payload.channel_id, new Set())
    }

    typingUsers.value.get(payload.channel_id)!.add(username)
  }

  function handleTypingStop(payload: TypingPayload) {
    if (!payload.channel_id || !payload.user_id || payload.user_id === authStore.user?.id) {
      return
    }

    const username = payload.username || payload.user?.username || 'Unknown User'
    const channelTypingUsers = typingUsers.value.get(payload.channel_id)

    if (channelTypingUsers) {
      channelTypingUsers.delete(username)
      if (channelTypingUsers.size === 0) {
        typingUsers.value.delete(payload.channel_id)
      }
    }
  }

  function getTypingUsers(channelId: string): string[] {
    const channelTypingUsers = typingUsers.value.get(channelId)
    return channelTypingUsers ? Array.from(channelTypingUsers) : []
  }

  function clearTypingUsers(channelId: string) {
    typingUsers.value.delete(channelId)
  }

  function initializeWebSocket() {
    websocketService.on('new_message', (payload: unknown) => {
      if (payload && typeof payload === 'object' && 'id' in payload) {
        addMessageFromWebSocket(payload as IncomingMessagePayload)
      }
    })

    websocketService.on('reaction_added', (payload: unknown) => {
      console.log('Received reaction_added event:', payload)
      if (payload && typeof payload === 'object' && 'message_id' in payload) {
        const reactionPayload = payload as IncomingReactionPayload
        addReactionToMessage(reactionPayload.message_id, reactionPayload)
      }
    })

    websocketService.on('reaction_removed', (payload: unknown) => {
      console.log('Received reaction_removed event:', payload)
      if (
        payload &&
        typeof payload === 'object' &&
        'message_id' in payload &&
        'user_id' in payload &&
        'reaction_type' in payload
      ) {
        const reactionPayload = payload as IncomingReactionPayload
        removeReactionFromMessage(
          reactionPayload.message_id,
          reactionPayload.user_id,
          reactionPayload.reaction_type,
        )
      } else {
        console.log('Invalid reaction_removed payload:', payload)
      }
    })

    websocketService.on('typing_start', (payload: unknown) => {
      console.log('Received typing_start event:', payload)
      if (
        payload &&
        typeof payload === 'object' &&
        'channel_id' in payload &&
        'user_id' in payload
      ) {
        handleTypingStart(payload as TypingPayload)
      }
    })

    websocketService.on('typing_stop', (payload: unknown) => {
      console.log('Received typing_stop event:', payload)
      if (
        payload &&
        typeof payload === 'object' &&
        'channel_id' in payload &&
        'user_id' in payload
      ) {
        handleTypingStop(payload as TypingPayload)
      }
    })

    websocketService.on('connected', () => {
      console.log('WebSocket connected, ready for real-time messaging')
    })
    websocketService.on('disconnected', () => {
      console.log('WebSocket disconnected')
    })
  }

  function cleanupWebSocket() {
    websocketService.off('new_message', (payload: unknown) => {
      if (payload && typeof payload === 'object' && 'id' in payload) {
        addMessageFromWebSocket(payload as IncomingMessagePayload)
      }
    })
    websocketService.off('reaction_added', (payload: unknown) => {
      if (payload && typeof payload === 'object' && 'message_id' in payload) {
        const reactionPayload = payload as IncomingReactionPayload
        addReactionToMessage(reactionPayload.message_id, reactionPayload)
      }
    })
    websocketService.off('reaction_removed', (payload: unknown) => {
      if (
        payload &&
        typeof payload === 'object' &&
        'message_id' in payload &&
        'user_id' in payload &&
        'reaction_type' in payload
      ) {
        const reactionPayload = payload as IncomingReactionPayload
        removeReactionFromMessage(
          reactionPayload.message_id,
          reactionPayload.user_id,
          reactionPayload.reaction_type,
        )
      }
    })
    websocketService.off('typing_start', (payload: unknown) => {
      if (
        payload &&
        typeof payload === 'object' &&
        'channel_id' in payload &&
        'user_id' in payload
      ) {
        handleTypingStart(payload as TypingPayload)
      }
    })
    websocketService.off('typing_stop', (payload: unknown) => {
      if (
        payload &&
        typeof payload === 'object' &&
        'channel_id' in payload &&
        'user_id' in payload
      ) {
        handleTypingStop(payload as TypingPayload)
      }
    })
  }

  return {
    messages,
    loading,
    currentChannelId,
    currentMessages,
    typingUsers,
    fetchMessages,
    addMessage,
    addReactionToMessage,
    removeReactionFromMessage,
    addReaction,
    removeReaction,
    clearMessages,
    sendMessage,
    startTyping,
    stopTyping,
    getTypingUsers,
    clearTypingUsers,
    initializeWebSocket,
    cleanupWebSocket,
  }
})
