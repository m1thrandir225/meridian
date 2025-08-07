<script setup lang="ts">
import ChannelsSidebar from '@/components/ChannelsSidebar.vue'
import MembersSidebar from '@/components/MembersSidebar.vue'
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'

// Accept props for sidebar state
defineProps<{
  isChannelsSidebarOpen: boolean
  isMembersSidebarOpen: boolean
}>()

// Emit events for sidebar toggles
defineEmits<{
  toggleChannelsSidebar: []
  toggleMembersSidebar: []
}>()
</script>

<template>
  <div class="flex h-screen">
    <SidebarProvider :default-open="true" variant="sidebar" collapsible="none">
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
        <slot
          :is-channels-sidebar-open="isChannelsSidebarOpen"
          :is-members-sidebar-open="isMembersSidebarOpen"
          :toggle-channels-sidebar="() => $emit('toggleChannelsSidebar')"
          :toggle-members-sidebar="() => $emit('toggleMembersSidebar')"
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
