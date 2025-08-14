<script setup lang="ts">
import ChatLayout from '@/layouts/ChatLayout.vue'
import ChatArea from '@/components/ChatArea.vue'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessageStore } from '@/stores/message'
import { useChannelStore } from '@/stores/channel'

const route = useRoute()
const messageStore = useMessageStore()
const channelStore = useChannelStore()

// Sidebar state
const isChannelsSidebarOpen = ref(true)
const isMembersSidebarOpen = ref(true)

watch(
  () => route.params.id,
  (id) => {
    if (id) {
      messageStore.fetchMessages(id as string)
      channelStore.setCurrentChannel(id as string)
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
