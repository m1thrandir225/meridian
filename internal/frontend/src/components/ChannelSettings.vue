<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Settings, Archive, ArchiveRestore, AlertTriangle } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useChannelStore } from '@/stores/channel'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'
import type { Channel } from '@/types/models/channel'

interface Props {
  channel: Channel
}

const props = defineProps<Props>()
const channelStore = useChannelStore()
const authStore = useAuthStore()

const showSettings = ref(false)
const isArchiving = ref(false)
const isUnarchiving = ref(false)

const isChannelOwner = computed(() => {
  return authStore.user?.id === props.channel.creator_user_id
})

const handleArchive = async () => {
  if (!confirm('Are you sure you want to archive this channel? This action can be undone.')) {
    return
  }

  isArchiving.value = true
  try {
    await channelStore.archiveChannel(props.channel.id)
    toast.success('Channel archived successfully')
    showSettings.value = false
  } catch (error) {
    console.error('Failed to unarchive channel:', error)
    toast.error('Failed to archive channel')
  } finally {
    isArchiving.value = false
  }
}

const handleUnarchive = async () => {
  if (!confirm('Are you sure you want to unarchive this channel?')) {
    return
  }

  isUnarchiving.value = true
  try {
    await channelStore.unarchiveChannel(props.channel.id)
    toast.success('Channel unarchived successfully')
    showSettings.value = false
  } catch (error) {
    console.error('Failed to unarchive channel:', error)
    toast.error('Failed to unarchive channel')
  } finally {
    isUnarchiving.value = false
  }
}
</script>

<template>
  <Dialog v-model:open="showSettings">
    <DialogTrigger as-child>
      <Button variant="ghost" size="icon" class="h-8 w-8">
        <Settings class="h-4 w-4" />
      </Button>
    </DialogTrigger>
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Channel Settings</DialogTitle>
        <DialogDescription> Manage settings for {{ channel.name }} </DialogDescription>
      </DialogHeader>

      <div class="space-y-4">
        <!-- Channel Info -->
        <div class="space-y-2">
          <h4 class="font-medium">Channel Information</h4>
          <div class="text-sm text-muted-foreground space-y-1">
            <div><strong>Name:</strong> {{ channel.name }}</div>
            <div><strong>Topic:</strong> {{ channel.topic || 'No topic set' }}</div>
            <div><strong>Members:</strong> {{ channel.members_count }}</div>
            <div>
              <strong>Created:</strong> {{ new Date(channel.creation_time).toLocaleDateString() }}
            </div>
            <div class="flex items-center gap-2">
              <strong>Status:</strong>
              <Badge :variant="channel.is_archived ? 'destructive' : 'secondary'">
                {{ channel.is_archived ? 'Archived' : 'Active' }}
              </Badge>
            </div>
          </div>
        </div>

        <Separator />

        <!-- Archive/Unarchive Actions -->
        <div v-if="isChannelOwner" class="space-y-2">
          <h4 class="font-medium">Channel Management</h4>

          <div v-if="!channel.is_archived" class="space-y-2">
            <div class="flex items-start gap-3 p-3 border rounded-lg bg-muted/50">
              <AlertTriangle class="h-5 w-5 text-amber-500 mt-0.5" />
              <div class="flex-1">
                <p class="text-sm font-medium">Archive Channel</p>
                <p class="text-xs text-muted-foreground">
                  Archived channels are hidden from the channel list but can be restored later.
                </p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              @click="handleArchive"
              :disabled="isArchiving"
              class="w-full"
            >
              <Archive v-if="!isArchiving" class="h-4 w-4 mr-2" />
              <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
              Archive Channel
            </Button>
          </div>

          <div v-else class="space-y-2">
            <div
              class="flex items-start gap-3 p-3 border rounded-lg bg-green-50 dark:bg-green-950/20"
            >
              <ArchiveRestore class="h-5 w-5 text-green-500 mt-0.5" />
              <div class="flex-1">
                <p class="text-sm font-medium">Restore Channel</p>
                <p class="text-xs text-muted-foreground">
                  This channel is currently archived. You can restore it to make it active again.
                </p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              @click="handleUnarchive"
              :disabled="isUnarchiving"
              class="w-full"
            >
              <ArchiveRestore v-if="!isUnarchiving" class="h-4 w-4 mr-2" />
              <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
              Restore Channel
            </Button>
          </div>
        </div>

        <div v-else class="text-sm text-muted-foreground">
          Only the channel owner can manage these settings.
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" @click="showSettings = false"> Close </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
