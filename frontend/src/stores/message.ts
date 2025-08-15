import channelService from '@/services/channel.service'
import websocketService from '@/services/websocket.service'
import type { Message } from '@/types/models/message'
import type { IncomingMessagePayload } from '@/types/websocket'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const loading = ref(false)
  const currentChannelId = ref<string | null>(null)

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
    const message: Message = {
      id: payload.id,
      channel_id: payload.channel_id,
      sender_user_id: payload.sender_id,
      content_text: payload.content,
      parent_message_id: payload.parent_message_id,
      created_at: payload.timestamp,
    }

    addMessage(message)
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

  function initializeWebSocket() {
    websocketService.on('new_message', (payload: unknown) => {
      if (payload && typeof payload === 'object' && 'id' in payload) {
        addMessageFromWebSocket(payload as IncomingMessagePayload)
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
  }

  return {
    messages,
    loading,
    currentChannelId,
    currentMessages,
    fetchMessages,
    addMessage,
    clearMessages,
    sendMessage,
    startTyping,
    stopTyping,
    initializeWebSocket,
    cleanupWebSocket,
  }
})
