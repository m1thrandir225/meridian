import channelService from '@/services/channel.service'
import type { Channel } from '@/types/models/channel'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChannelStore = defineStore('channels', () => {
  const channels = ref<Channel[]>([])
  const currentChannel = ref<Channel | null>(null)
  const loading = ref(false)

  function setCurrentChannel(id: string) {
    currentChannel.value = channels.value.filter((item) => item.id === id)[0]
  }

  async function fetchChannels() {
    loading.value = true
    try {
      const response = await channelService.getChannels()
      channels.value = response
    } catch (error) {
      console.error('Failed to fetch channels:', error)
    } finally {
      loading.value = false
    }
  }
  async function addChannel(newChannel: Channel) {
    channels.value.push(newChannel)
  }

  return {
    channels,
    loading,
    currentChannel,
    fetchChannels,
    addChannel,
    setCurrentChannel,
  }
})
