<script setup lang="ts">
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import integrationService from '@/services/integration.service'
import { useChannelStore } from '@/stores/channel'
import type {
  CreateIntegrationRequest,
  CreateIntegrationResponse,
} from '@/types/responses/integration'
import { useMutation } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import * as z from 'zod'
import { Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()

const newBotSchema = toTypedSchema(
  z.object({
    serviceName: z.string().min(5),
    targetChannelIds: z.array(z.string()).min(1),
  }),
)

const { handleSubmit, isFieldDirty, values, setFieldValue } = useForm({
  validationSchema: newBotSchema,
  initialValues: {
    serviceName: '',
    targetChannelIds: [],
  },
})

const createdIntegration = ref<CreateIntegrationResponse | null>(null)
const showSuccess = computed(() => createdIntegration.value !== null)

const { mutateAsync, status } = useMutation({
  mutationKey: ['create-bot'],
  mutationFn: (input: CreateIntegrationRequest) => integrationService.createIntegration(input),
  onSuccess: async (data: CreateIntegrationResponse) => {
    createdIntegration.value = data
    await nextTick()
    console.log('createdIntegration', createdIntegration.value)
    toast.success('Bot created successfully!')
  },
  onError: (error) => {
    toast.error(error.message || 'Failed to create bot')
    console.error('Error creating bot:', error)
  },
})

const onSubmit = handleSubmit(async (values) => {
  await mutateAsync({
    service_name: values.serviceName,
    target_channel_ids: values.targetChannelIds,
  })
})

const channelStore = useChannelStore()
const authStore = useAuthStore()
const channels = computed(() => channelStore.channels)
const isFetchingChannels = computed(() => channelStore.loading)
const ownerChannels = computed(() =>
  channels.value.filter((c) => c.creator_user_id === authStore.user?.id),
)

onMounted(async () => {
  await channelStore.fetchChannels()
})

const toggleChannelSelection = (channelId: string, checked: boolean | string) => {
  const isChecked = Boolean(checked)
  const currentIds = values.targetChannelIds || []
  let newIds: string[]

  if (isChecked) {
    newIds = [...currentIds, channelId]
  } else {
    newIds = currentIds.filter((id) => id !== channelId)
  }

  setFieldValue('targetChannelIds', newIds)
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

const createAnotherBot = () => {
  createdIntegration.value = null
  setFieldValue('serviceName', '')
  setFieldValue('targetChannelIds', [])
}

// Add this debugging
console.log('Component state:', {
  createdIntegration: createdIntegration.value,
  isFetchingChannels: isFetchingChannels.value,
  channels: channels.value?.length || 0,
  status: status.value,
})
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

      <!-- Show loading state -->
      <Card v-if="isFetchingChannels">
        <CardContent class="p-6">
          <div class="flex items-center justify-center">
            <Loader2 class="animate-spin mr-2" />
            <span>Loading channels...</span>
          </div>
        </CardContent>
      </Card>

      <!-- Show form -->
      <Card v-else-if="!showSuccess">
        <CardHeader>
          <CardTitle>Create New Bot</CardTitle>
          <CardDescription>
            Configure your integration bot by providing a service name and selecting target
            channels.
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <form @submit="onSubmit" class="space-y-6" v-if="!isFetchingChannels">
            <div class="space-y-2">
              <FormField
                v-slot="{ componentField }"
                name="serviceName"
                :validate-on-blur="!isFieldDirty"
              >
                <FormItem>
                  <FormLabel>Service Name</FormLabel>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="e.g., GitHub Bot, Slack Bot, Custom API"
                      v-bind="componentField"
                    />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <div class="space-y-2" v-if="channels && channels.length > 0">
              <FormField name="targetChannelIds">
                <FormItem>
                  <div class="mb-4">
                    <FormLabel class="text-base"> Target Channels </FormLabel>
                    <FormDescription>
                      Select the channels you want to send messages to.
                    </FormDescription>
                  </div>
                  <div class="space-y-3">
                    <div
                      v-for="channel in ownerChannels"
                      :key="channel.id"
                      class="flex flex-row items-start space-x-3"
                    >
                      <FormControl>
                        <Checkbox
                          :model-value="values.targetChannelIds?.includes(channel.id) || false"
                          @update:model-value="
                            (checked) => toggleChannelSelection(channel.id, checked)
                          "
                        />
                      </FormControl>
                      <FormLabel class="font-normal">
                        {{ channel.name }}
                      </FormLabel>
                    </div>
                  </div>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <Button type="submit" :disabled="status === 'pending'" class="w-full">
              <Loader2 v-if="status === 'pending'" class="animate-spin" />
              <span v-else>Create Bot</span>
            </Button>
          </form>
        </CardContent>
      </Card>

      <!-- Show success state -->
      <Card v-else>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <span>Bot Created Successfully!</span>
            <Badge variant="secondary">{{ createdIntegration?.service_name }}</Badge>
          </CardTitle>
          <CardDescription>
            Your integration bot has been created. Use the token below to send messages via API.
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-2">
            <Label>Integration Token</Label>
            <div class="flex gap-2">
              <Input
                :default-value="createdIntegration?.token"
                readonly
                class="font-mono text-sm"
              />
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
                v-for="channelId in createdIntegration?.target_channels"
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
  -H "Authorization: ApiKey {{ createdIntegration?.token }}" \
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
  -H "Authorization: ApiKey {{ createdIntegration?.token }}" \
  -d '{
    "content_text": "Hello from your bot!",
    "target_channel_id": "{{ createdIntegration?.target_channels[0] }}"
  }'</code></pre>
                </div>
              </div>
            </div>
          </div>

          <div class="flex gap-2">
            <Button @click="goToChat" class="flex-1"> Go to Chat </Button>
            <Button variant="outline" @click="createAnotherBot"> Create Another Bot </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
