<script setup lang="ts">
import ChannelsSidebar from '@/components/ChannelsSidebar.vue'
import MembersSidebar from '@/components/MembersSidebar.vue'
import ChatArea from '@/components/ChatArea.vue'
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'
import { ref } from 'vue'

const isChannelsSidebarOpen = ref(true)
const isMembersSidebarOpen = ref(true)

const toggleChannelsSidebar = () => {
  isChannelsSidebarOpen.value = !isChannelsSidebarOpen.value
}

const toggleMembersSidebar = () => {
  isMembersSidebarOpen.value = !isMembersSidebarOpen.value
}
</script>

<template>
  <div class="flex h-screen">
    <SidebarProvider>
      <Transition
        name="slide-left"
        enter-active-class="transition-all duration-300 ease-out"
        leave-active-class="transition-all duration-300 ease-in"
        enter-from-class="-translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="-translate-x-full opacity-0"
      >
        <ChannelsSidebar v-if="isChannelsSidebarOpen" />
      </Transition>
      <SidebarInset>
        <ChatArea
          :is-channels-sidebar-open="isChannelsSidebarOpen"
          :is-members-sidebar-open="isMembersSidebarOpen"
          @toggle-channels-sidebar="toggleChannelsSidebar"
          @toggle-members-sidebar="toggleMembersSidebar"
        />
      </SidebarInset>
      <Transition
        name="slide-right"
        enter-active-class="transition-all duration-300 ease-out"
        leave-active-class="transition-all duration-300 ease-in"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
      >
        <MembersSidebar v-if="isMembersSidebarOpen" />
      </Transition>
    </SidebarProvider>
  </div>
</template>
