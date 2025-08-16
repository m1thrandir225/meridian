<script setup lang="ts">
import channelService from '@/services/channel.service'
import type { Channel } from '@/types/models/channel'
import type { CreateIntegrationResponse } from '@/types/responses/integration'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'

// Placeholder for bot registration page
const router = useRouter()

const serviceName = ref('')
const selectedChannels = ref<string[]>([])
const channels = ref<Channel[]>([])
const isLoading = ref(false)
const createdIntegration = ref<CreateIntegrationResponse | null>(null)

onMounted(async () => {
  try {
    const userChannels = await channelService.getChannels()
    channels.value = userChannels
  } catch (error) {
    toast.error('Failed to load channels')
    console.error('Error loading channels:', error)
  }
})

const toggleChannelSelection = (channelId: string) => {
  const index = selectedChannels.value.indexOf(channelId)
  if (index > -1) {
    selectedChannels.value.splice(index, 1)
  } else {
    selectedChannels.value.push(channelId)
  }
}

const createBot = async () => {
  if (!serviceName.value.trim()) {
    toast.error('Please enter a service name')
    return
  }

  if (selectedChannels.value.length === 0) {
    toast.error('Please select at least one channel')
    return
  }

  isLoading.value = true

  // try {
  //   const request: CreateIntegrationRequest = {
  //     service_name: serviceName.value.trim(),
  //     target_channel_ids: selectedChannels.value,
  //   }

  //   const response = await integrationService.createIntegration(request)
  //   createdIntegration.value = response
  //   toast.success('Bot created successfully!')
  // } catch (error: any) {
  //   toast.error(error.response?.data?.error || 'Failed to create bot')
  //   console.error('Error creating bot:', error)
  // } finally {
  //   isLoading.value = false
  // }
}

const copyToken = async () => {
  if (createdIntegration.value?.token) {
    try {
      await navigator.clipboard.writeText(createdIntegration.value.token)
      toast.success('Token copied to clipboard!')
    } catch (error) {
      toast.error('Failed to copy token: ' + error)
    }
  }
}

const goToChat = () => {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-background flex items-center justify-center p-4">
    <div class="max-w-2xl w-full mx-auto">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold mb-2">Bot Registration</h1>
        <p class="text-muted-foreground">
          Create a new integration bot to send messages to your channels
        </p>
      </div>

      <Card v-if="!createdIntegration">
        <CardHeader>
          <CardTitle>Create New Bot</CardTitle>
          <CardDescription>
            Configure your integration bot by providing a service name and selecting target
            channels.
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-2">
            <Label for="service-name">Service Name</Label>
            <Input
              id="service-name"
              v-model="serviceName"
              placeholder="e.g., GitHub Bot, Slack Bot, Custom API"
              :disabled="isLoading"
            />
          </div>

          <div class="space-y-2">
            <Label>Target Channels</Label>
            <div class="grid gap-2 max-h-60 overflow-y-auto border rounded-md p-4">
              <div
                v-for="channel in channels"
                :key="channel.id"
                class="flex items-center space-x-2 p-2 rounded-md hover:bg-muted cursor-pointer"
                @click="toggleChannelSelection(channel.id)"
              >
                <input
                  type="checkbox"
                  :checked="selectedChannels.includes(channel.id)"
                  class="rounded"
                  @change="toggleChannelSelection(channel.id)"
                />
                <div class="flex-1">
                  <div class="font-medium">{{ channel.name }}</div>
                  <div class="text-sm text-muted-foreground">{{ channel.topic || 'No topic' }}</div>
                </div>
              </div>
            </div>
            <p class="text-sm text-muted-foreground">
              Selected: {{ selectedChannels.length }} channel{{
                selectedChannels.length !== 1 ? 's' : ''
              }}
            </p>
          </div>

          <Button
            @click="createBot"
            :disabled="isLoading || !serviceName.trim() || selectedChannels.length === 0"
            class="w-full"
          >
            <span v-if="isLoading">Creating...</span>
            <span v-else>Create Bot</span>
          </Button>
        </CardContent>
      </Card>

      <Card v-else>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <span>Bot Created Successfully!</span>
            <Badge variant="secondary">{{ createdIntegration.service_name }}</Badge>
          </CardTitle>
          <CardDescription>
            Your integration bot has been created. Use the token below to send messages via API.
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-2">
            <Label>Integration Token</Label>
            <div class="flex gap-2">
              <Input :value="createdIntegration.token" readonly class="font-mono text-sm" />
              <Button variant="outline" @click="copyToken"> Copy </Button>
            </div>
            <p class="text-sm text-muted-foreground">
              Keep this token secure. You'll need it to send messages via the API.
            </p>
          </div>

          <Separator />

          <div class="space-y-2">
            <Label>Target Channels</Label>
            <div class="flex flex-wrap gap-2">
              <Badge
                v-for="channelId in createdIntegration.target_channel_ids"
                :key="channelId"
                variant="outline"
              >
                {{ channels.find((c) => c.id === channelId)?.name || channelId }}
              </Badge>
            </div>
          </div>

          <div class="space-y-4">
            <h4 class="font-medium">API Usage Examples</h4>
            <div class="space-y-4">
              <div>
                <h5 class="text-sm font-medium mb-2">Webhook Message:</h5>
                <div class="bg-muted p-4 rounded-md">
                  <pre
                    class="text-sm overflow-x-auto"
                  ><code>curl -X POST http://api.localhost/api/v1/integrations/webhook/message \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey {{ createdIntegration.token }}" \
  -d '{
    "content_text": "Hello from your bot!"
  }'</code></pre>
                </div>
              </div>

              <div>
                <h5 class="text-sm font-medium mb-2">Callback Message:</h5>
                <div class="bg-muted p-4 rounded-md">
                  <pre
                    class="text-sm overflow-x-auto"
                  ><code>curl -X POST http://api.localhost/api/v1/integrations/callback/message \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey {{ createdIntegration.token }}" \
  -d '{
    "content_text": "Hello from your bot!",
    "target_channel_id": "{{ createdIntegration.target_channel_ids[0] }}"
  }'</code></pre>
                </div>
              </div>
            </div>
          </div>

          <div class="flex gap-2">
            <Button @click="goToChat" class="flex-1"> Go to Chat </Button>
            <Button variant="outline" @click="createdIntegration = null">
              Create Another Bot
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
