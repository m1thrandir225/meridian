<script setup lang="ts">
import ChatLayout from '@/layouts/ChatLayout.vue'
import ChatArea from '@/components/ChatArea.vue'
import { ref } from 'vue'

// Sidebar state
const isChannelsSidebarOpen = ref(true)
const isMembersSidebarOpen = ref(true)
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
