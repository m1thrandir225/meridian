<script setup lang="ts">
import { computed } from 'vue'
import { useMessageStore } from '@/stores/message'

interface Props {
  channelId: string
}

const props = defineProps<Props>()
const messageStore = useMessageStore()

const typingUsers = computed(() => messageStore.getTypingUsers(props.channelId))
const isVisible = computed(() => typingUsers.value.length > 0)

const typingText = computed(() => {
  const users = typingUsers.value
  if (users.length === 0) return ''

  if (users.length === 1) {
    return `${users[0]} is typing...`
  } else if (users.length === 2) {
    return `${users[0]} and ${users[1]} are typing...`
  } else {
    return `${users[0]} and ${users.length - 1} others are typing...`
  }
})
</script>

<template>
  <Transition
    enter-active-class="transition-all duration-300 ease-out"
    enter-from-class="opacity-0 transform translate-y-2"
    enter-to-class="opacity-100 transform translate-y-0"
    leave-active-class="transition-all duration-200 ease-in"
    leave-from-class="opacity-100 transform translate-y-0"
    leave-to-class="opacity-0 transform translate-y-2"
  >
    <div
      v-if="isVisible"
      class="flex items-center gap-2 px-4 py-2 text-sm text-muted-foreground bg-muted/30 border-t border-border"
    >
      <div class="flex items-center gap-1">
        <div class="flex space-x-1">
          <div
            class="w-2 h-2 bg-muted-foreground rounded-full animate-bounce"
            style="animation-delay: 0ms"
          ></div>
          <div
            class="w-2 h-2 bg-muted-foreground rounded-full animate-bounce"
            style="animation-delay: 150ms"
          ></div>
          <div
            class="w-2 h-2 bg-muted-foreground rounded-full animate-bounce"
            style="animation-delay: 300ms"
          ></div>
        </div>
      </div>
      <span class="text-xs">{{ typingText }}</span>
    </div>
  </Transition>
</template>
