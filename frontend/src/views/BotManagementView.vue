<script setup lang="ts">
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import integrationService from '@/services/integration.service'
import { useChannelStore } from '@/stores/channel'
import type { IntegrationBot } from '@/types/models/integration_bot'
import type { UpdateIntegrationRequest } from '@/types/responses/integration'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { AlertTriangle, Edit, Loader2, Trash2 } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'

const queryClient = useQueryClient()
const channelStore = useChannelStore()

// State
const editingIntegration = ref<IntegrationBot | null>(null)
const showEditDialog = ref(false)
const selectedChannels = ref<string[]>([])

// Queries
const {
  data: integrationsData,
  isLoading,
  error,
} = useQuery({
  queryKey: ['integrations'],
  queryFn: () => integrationService.listIntegrations(),
})

const integrations = computed(() => integrationsData.value?.integrations || [])

// Mutations
const revokeMutation = useMutation({
  mutationFn: (integrationId: string) =>
    integrationService.revokeIntegration({ integration_id: integrationId }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['integrations'] })
    toast.success('Bot revoked successfully')
  },
  onError: (error) => {
    toast.error('Failed to revoke bot: ' + error.message)
  },
})

const updateMutation = useMutation({
  mutationFn: (data: UpdateIntegrationRequest) => integrationService.updateIntegration(data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['integrations'] })
    showEditDialog.value = false
    editingIntegration.value = null
    toast.success('Bot updated successfully')
  },
  onError: (error) => {
    toast.error('Failed to update bot: ' + error.message)
  },
})

// Computed
const channels = computed(() => channelStore.channels)
const isFetchingChannels = computed(() => channelStore.loading)

// Methods
const handleRevoke = (integration: IntegrationBot) => {
  if (
    confirm(
      `Are you sure you want to revoke the "${integration.service_name}" bot? This action cannot be undone.`,
    )
  ) {
    revokeMutation.mutate(integration.id)
  }
}

const handleEdit = (integration: IntegrationBot) => {
  editingIntegration.value = integration
  selectedChannels.value = [...integration.target_channels]
  showEditDialog.value = true
}

const handleUpdate = () => {
  if (!editingIntegration.value) return

  updateMutation.mutate({
    integration_id: editingIntegration.value.id,
    target_channel_ids: selectedChannels.value,
  })
}

const toggleChannelSelection = (channelId: string, checked: boolean | string) => {
  if (checked) {
    selectedChannels.value = [...selectedChannels.value, channelId]
  } else {
    selectedChannels.value = selectedChannels.value.filter((id) => id !== channelId)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Lifecycle
onMounted(async () => {
  await channelStore.fetchChannels()
})
</script>

<template>
  <div class="min-h-screen bg-background p-6">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold mb-2">Bot Management</h1>
        <p class="text-muted-foreground">
          Manage your integration bots, update target channels, and revoke access when needed.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center py-12">
        <Loader2 class="animate-spin mr-2" />
        <span>Loading bots...</span>
      </div>

      <!-- Error State -->
      <Alert v-else-if="error" variant="destructive" class="mb-6">
        <AlertTriangle class="h-4 w-4" />
        <AlertDescription> Failed to load bots: {{ error.message }} </AlertDescription>
      </Alert>

      <!-- Empty State -->
      <Card v-else-if="integrations.length === 0">
        <CardContent class="p-12 text-center">
          <div class="mb-4">
            <div
              class="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <AlertTriangle class="w-8 h-8 text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold mb-2">No bots found</h3>
            <p class="text-muted-foreground mb-6">You haven't created any integration bots yet.</p>
            <Button @click="$router.push('/bot-registration')"> Create Your First Bot </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Bots List -->
      <div v-else class="space-y-6">
        <div
          v-for="integration in integrations"
          :key="integration.id"
          class="border rounded-lg p-6 bg-card"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h3 class="text-xl font-semibold">{{ integration.service_name }}</h3>
                <Badge :variant="integration.is_revoked ? 'destructive' : 'secondary'">
                  {{ integration.is_revoked ? 'Revoked' : 'Active' }}
                </Badge>
              </div>
              <p class="text-sm text-muted-foreground">
                Created {{ formatDate(integration.created_at) }}
              </p>
            </div>

            <div class="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                @click="handleEdit(integration)"
                :disabled="integration.is_revoked"
              >
                <Edit class="w-4 h-4 mr-1" />
                Edit
              </Button>
              <Button
                variant="destructive"
                size="sm"
                @click="handleRevoke(integration)"
                :disabled="integration.is_revoked"
              >
                <Trash2 class="w-4 h-4 mr-1" />
                Revoke
              </Button>
            </div>
          </div>

          <Separator class="my-4" />

          <!-- Target Channels -->
          <div>
            <h4 class="font-medium mb-3">Target Channels</h4>
            <div class="flex flex-wrap gap-2">
              <Badge
                v-for="channelId in integration.target_channels"
                :key="channelId"
                variant="outline"
              >
                {{ channels.find((c) => c.id === channelId)?.name || channelId }}
              </Badge>
            </div>
          </div>
        </div>
      </div>

      <!-- Edit Dialog -->
      <Dialog v-model:open="showEditDialog">
        <DialogContent class="max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Bot Channels</DialogTitle>
            <DialogDescription>
              Update the target channels for {{ editingIntegration?.service_name }}
            </DialogDescription>
          </DialogHeader>

          <div class="space-y-4">
            <div v-if="isFetchingChannels" class="flex items-center justify-center py-4">
              <Loader2 class="animate-spin mr-2" />
              <span>Loading channels...</span>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="channel in channels"
                :key="channel.id"
                class="flex items-center space-x-3"
              >
                <Checkbox
                  :model-value="selectedChannels.includes(channel.id)"
                  @update:model-value="(checked) => toggleChannelSelection(channel.id, checked)"
                  :id="channel.id"
                />
                <Label :for="channel.id" class="font-normal">
                  {{ channel.name }}
                </Label>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              @click="showEditDialog = false"
              :disabled="updateMutation.isPending"
            >
              Cancel
            </Button>
            <Button
              @click="handleUpdate"
              :disabled="updateMutation.isPending || selectedChannels.length === 0"
            >
              <Loader2 v-if="updateMutation.isPending" class="animate-spin mr-2" />
              Update Bot
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  </div>
</template>
