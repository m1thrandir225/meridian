import channelService from '@/services/channel.service'
import type { Channel } from '@/types/models/channel'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useChannelStore = defineStore('channels', () => {
  const channels = ref<Channel[]>([])
  const currentChannel = ref<Channel | null>(null)
  const loading = ref(false)

  const getCurrentChannel = computed(() => currentChannel.value)

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
  function addChannel(newChannel: Channel) {
    channels.value.push(newChannel)
  }

  function updateChannel(updatedChannel: Channel) {
    const index = channels.value.findIndex((c) => c.id === updatedChannel.id)
    if (index !== -1) {
      channels.value[index] = updatedChannel
      if (currentChannel.value?.id === updatedChannel.id) {
        currentChannel.value = updatedChannel
      }
    }
  }
  function removeChannel(channelId: string) {
    channels.value = channels.value.filter((c) => c.id !== channelId)
    if (currentChannel.value?.id === channelId) {
      currentChannel.value = null
    }
  }
  async function archiveChannel(channelId: string) {
    try {
      await channelService.archiveChannel(channelId)
      const channel = channels.value.find((c) => c.id === channelId)
      if (channel) {
        channel.is_archived = true
        updateChannel(channel)
      }
    } catch (error) {
      console.error('Failed to archive channel:', error)
      throw error
    }
  }

  async function unarchiveChannel(channelId: string) {
    try {
      await channelService.unarchiveChannel(channelId)
      // Update the channel in the store
      const channel = channels.value.find((c) => c.id === channelId)
      if (channel) {
        channel.is_archived = false
        updateChannel(channel)
      }
    } catch (error) {
      console.error('Failed to unarchive channel:', error)
      throw error
    }
  }

  return {
    channels,
    currentChannel,
    loading,
    getCurrentChannel,
    setCurrentChannel,
    fetchChannels,
    addChannel,
    updateChannel,
    removeChannel,
    archiveChannel,
    unarchiveChannel,
  }
})
