import channelService from '@/services/channel.service'
import type { Message } from '@/types/models/message'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const loading = ref(false)

  async function fetchMessages(channelId: string) {
    loading.value = true
    messages.value = await channelService.getMessages(channelId)
    loading.value = false
  }

  return { messages, loading, fetchMessages }
})
