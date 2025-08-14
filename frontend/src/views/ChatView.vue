<script setup lang="ts">
import ChatLayout from '@/layouts/ChatLayout.vue'
import ChatArea from '@/components/ChatArea.vue'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessageStore } from '@/stores/message'
import { useChannelStore } from '@/stores/channel'

const route = useRoute()
const messageStore = useMessageStore()
const channelStore = useChannelStore()

// Sidebar state
const isChannelsSidebarOpen = ref(true)
const isMembersSidebarOpen = ref(true)

onMounted(async () => {
  if (!channelStore.channels.length) {
    await channelStore.fetchChannels()
  }
  if (route.params.id) {
    channelStore.setCurrentChannel(route.params.id as string)
    messageStore.fetchMessages(route.params.id as string)
  }
})
watch(
  [() => route.params.id, () => channelStore.channels.length],
  ([id, channelsLength], [oldId, oldChannelsLength]) => {
    if (id && channelsLength > 0) {
      channelStore.setCurrentChannel(id as string)
      messageStore.fetchMessages(id as string)
    }
  },
  { immediate: true },
)
</script>

<template>
  <ChatLayout
    :is-channels-sidebar-open="isChannelsSidebarOpen"
    :is-members-sidebar-open="isMembersSidebarOpen"
    @toggle-channels-sidebar="isChannelsSidebarOpen = !isChannelsSidebarOpen"
    @toggle-members-sidebar="isMembersSidebarOpen = !isMembersSidebarOpen"
  >
    <template #default="{ toggleChannelsSidebar, toggleMembersSidebar }">
      <ChatArea
        :is-channels-sidebar-open="isChannelsSidebarOpen"
        :is-members-sidebar-open="isMembersSidebarOpen"
        @toggle-channels-sidebar="toggleChannelsSidebar"
        @toggle-members-sidebar="toggleMembersSidebar"
      />
    </template>
  </ChatLayout>
</template>
