<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'
import { Reply } from 'lucide-vue-next'
import { useMessageStore } from '@/stores/message'
import type { Message } from '@/types/models/message'

interface Props {
  message: Message
}

const props = defineProps<Props>()
const emit = defineEmits<{
  reply: [message: Message]
}>()

const messageStore = useMessageStore()

const handleReactionClick = (reactionType: string) => {
  messageStore.addReaction(props.message.id, reactionType)
}

const quickReactions = ['👍', '👎', '❤️', '😂', '😮', '😢', '😡', '🤔', '👏', '👀']

const handleReply = () => {
  emit('reply', props.message)
}
</script>

<template>
  <HoverCard>
    <HoverCardTrigger as-child>
      <slot />
    </HoverCardTrigger>
    <HoverCardContent class="w-64 p-3" side="top" align="end">
      <!-- Quick Reactions -->
      <div class="mb-3">
        <div class="text-xs font-medium text-muted-foreground mb-2">Quick Reactions</div>
        <div class="flex flex-wrap gap-1">
          <Button
            v-for="emoji in quickReactions"
            :key="emoji"
            variant="ghost"
            size="sm"
            class="h-8 w-8 p-0 text-sm hover:scale-110 transition-transform duration-150"
            @click="handleReactionClick(emoji)"
          >
            {{ emoji }}
          </Button>
        </div>
      </div>

      <!-- Reply Action -->
      <div class="pt-2 border-t">
        <Button
          variant="ghost"
          size="sm"
          class="w-full justify-start h-8 text-sm"
          @click="handleReply"
        >
          <Reply class="h-4 w-4 mr-2" />
          Reply to message
        </Button>
      </div>
    </HoverCardContent>
  </HoverCard>
</template>
