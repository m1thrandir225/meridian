<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import inviteService from '@/services/invites.service'
import { useChannelStore } from '@/stores/channel'
import type { AcceptChannelInviteRequest } from '@/types/responses/channel_invite'
import { useMutation } from '@tanstack/vue-query'
import { CheckCircle, Loader2, XCircle } from 'lucide-vue-next'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'

//Stores & Router
const route = useRoute()
const router = useRouter()
const channelStore = useChannelStore()

const { mutateAsync, status, error } = useMutation({
  mutationKey: ['accept-invite'],
  mutationFn: (input: AcceptChannelInviteRequest) => inviteService.acceptInvite(input),
  onSuccess: async (data) => {
    toast.success('Invite accepted')
    await channelStore.fetchChannels()
    channelStore.setCurrentChannel(data.id)
    setTimeout(() => {
      router.push({ name: 'channel', params: { id: data.id } })
    }, 2000)
  },
  onError: (error) => {
    toast.error('Failed to accept invite')
    console.error(error)
  },
})

onMounted(async () => {
  const inviteCode = route.params.inviteCode as string

  if (!inviteCode) {
    toast.error('Invalid invite link')
    return
  }
  await mutateAsync({ invite_code: inviteCode })
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle>Channel Invite</CardTitle>
        <CardDescription>
          {{
            status === 'pending'
              ? 'Processing your invite...'
              : status === 'success'
                ? 'Successfully joined the channel!'
                : 'Unable to join the channel'
          }}
        </CardDescription>
      </CardHeader>
      <CardContent class="text-center space-y-4">
        <!-- Loading State -->
        <div v-if="status === 'pending'" class="flex flex-col items-center space-y-2">
          <Loader2 class="h-8 w-8 animate-spin text-primary" />
          <p class="text-sm text-muted-foreground">Accepting invite...</p>
        </div>

        <!-- Success State -->
        <div v-else-if="status === 'success'" class="flex flex-col items-center space-y-2">
          <CheckCircle class="h-8 w-8 text-green-500" />
          <p class="text-sm text-muted-foreground">Redirecting to channel...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="status === 'error'" class="flex flex-col items-center space-y-2">
          <XCircle class="h-8 w-8 text-red-500" />
          <p class="text-sm text-red-600">{{ error?.message || 'Failed to accept invite' }}</p>
          <Button @click="router.push('/')" variant="outline"> Go to Home </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
