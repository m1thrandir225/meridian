<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth'
import { useMessageStore } from '@/stores/message'
import type { Message } from '@/types/models/message'

interface Props {
  message: Message
}

const props = defineProps<Props>()
const authStore = useAuthStore()
const messageStore = useMessageStore()

const currentUserId = computed(() => authStore.user?.id)

const reactionGroups = computed(() => {
  if (!props.message.reactions || props.message.reactions.length === 0) {
    return []
  }

  const groups = new Map<string, { count: number; users: string[] }>()

  props.message.reactions.forEach((reaction) => {
    if (!groups.has(reaction.reaction_type)) {
      groups.set(reaction.reaction_type, { count: 0, users: [] })
    }
    const group = groups.get(reaction.reaction_type)!
    group.count++
    group.users.push(reaction.user_id)
  })

  return Array.from(groups.entries()).map(([type, data]) => ({
    type,
    count: data.count,
    hasUserReacted: data.users.includes(currentUserId.value || ''),
  }))
})

const handleReactionClick = (reactionType: string) => {
  const hasReacted = reactionGroups.value.find((g) => g.type === reactionType)?.hasUserReacted

  if (hasReacted) {
    messageStore.removeReaction(props.message.id, reactionType)
  } else {
    messageStore.addReaction(props.message.id, reactionType)
  }
}

const quickReactions = ['👍', '❤️', '😂', '😮', '😢', '😡']
</script>

<template>
  <div class="flex flex-wrap gap-1 mt-2">
    <!-- Existing reactions -->
    <Button
      v-for="group in reactionGroups"
      :key="group.type"
      variant="secondary"
      size="sm"
      class="h-6 px-2 text-xs"
      :class="{ 'bg-primary text-primary-foreground': group.hasUserReacted }"
      @click="handleReactionClick(group.type)"
    >
      {{ group.type }} {{ group.count }}
    </Button>

    <!-- Quick reaction buttons (visible on hover) -->
    <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
      <Button
        v-for="emoji in quickReactions"
        :key="emoji"
        variant="ghost"
        size="sm"
        class="h-6 w-6 p-0 text-xs"
        @click="handleReactionClick(emoji)"
      >
        {{ emoji }}
      </Button>
    </div>
  </div>
</template>
