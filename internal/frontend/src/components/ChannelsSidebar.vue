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
import { Input } from '@/components/ui/input'
import type { SidebarProps } from '@/components/ui/sidebar'
import { Bot, Loader2, Plus, Settings } from 'lucide-vue-next'
import { onMounted, ref, watch } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader } from '@/components/ui/sidebar'
import channelService from '@/services/channel.service'
import type { CreateChannelRequest } from '@/types/responses/channel'
import { useMutation } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'
import NavUser from './NavUser.vue'
import { FormControl, FormField, FormItem, FormLabel } from './ui/form'
import { toast } from 'vue-sonner'
import { useChannelStore } from '@/stores/channel'
import Logo from './LogoWord.vue'

const props = withDefaults(
  defineProps<
    SidebarProps & {
      showCreateDialog?: boolean
    }
  >(),
  {
    collapsible: 'icon',
    side: 'left',
    showCreateDialog: false,
  },
)

const emit = defineEmits<{
  'update:showCreateDialog': [value: boolean]
}>()

const router = useRouter()
const channelStore = useChannelStore()

const isNewChannelDialogOpen = ref(false)

watch(
  () => props.showCreateDialog,
  (newValue) => {
    if (newValue) {
      isNewChannelDialogOpen.value = true
    }
  },
)

watch(
  () => isNewChannelDialogOpen.value,
  (newValue) => {
    emit('update:showCreateDialog', newValue)
  },
)

onMounted(() => channelStore.fetchChannels())

const createChannelSchema = toTypedSchema(
  z.object({
    name: z.string().min(5),
    topic: z.string().min(5),
  }),
)

const { isFieldDirty, handleSubmit } = useForm({
  validationSchema: createChannelSchema,
})

const { mutateAsync, status } = useMutation({
  mutationKey: ['createChannel'],
  mutationFn: (input: CreateChannelRequest) => channelService.createChannel(input),
  onSuccess: async (response) => {
    isNewChannelDialogOpen.value = false
    toast.success('Sucessfully created a new channel')

    await channelStore.fetchChannels()

    setTimeout(() => {
      router.push({
        name: 'channel',
        params: { id: response.id },
      })
    }, 1000)
  },
  onError: (error) => {
    toast.error(error.message)
  },
})

const createChannel = handleSubmit(async (values) => {
  const { name, topic } = values
  await mutateAsync({ name, topic })
})
</script>

<template>
  <Sidebar v-bind="props" class="border-r">
    <SidebarHeader>
      <div class="flex items-center justify-between gap-3 px-4 py-3 border-b">
        <Logo size="24" />
        <Button variant="ghost" size="icon" class="h-6 w-6">
          <Settings class="h-4 w-4" />
        </Button>
      </div>
    </SidebarHeader>

    <SidebarContent class="flex flex-col">
      <div class="p-2 flex-1 flex flex-col min-h-0">
        <div class="flex-1 flex flex-col min-h-0">
          <div class="flex items-center justify-between px-2 py-1 mb-2 flex-shrink-0">
            <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Channels
            </h3>
            <Dialog v-model:open="isNewChannelDialogOpen">
              <DialogTrigger as-child>
                <Button variant="ghost" size="icon" class="h-4 w-4">
                  <Plus class="h-3 w-3" />
                </Button>
              </DialogTrigger>
              <DialogContent class="sm:max-w-[425px]">
                <DialogHeader>
                  <DialogTitle>Create New Channel</DialogTitle>
                  <DialogDescription>
                    Add a new text channel to your server. Channel names should be lowercase and
                    descriptive.
                  </DialogDescription>
                </DialogHeader>
                <form id="channelForm" class="grid gap-4 py-4" @submit="createChannel">
                  <FormField
                    v-slot="{ componentField }"
                    name="name"
                    :validate-on-blur="!isFieldDirty"
                  >
                    <FormItem>
                      <FormLabel for="password">Name</FormLabel>
                      <FormControl>
                        <Input type="text" placeholder="meridian-chats" v-bind="componentField" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                  <FormField
                    v-slot="{ componentField }"
                    name="topic"
                    :validate-on-blur="!isFieldDirty"
                  >
                    <FormItem>
                      <FormLabel for="topic">Topic</FormLabel>
                      <FormControl>
                        <Input type="text" placeholder="general" v-bind="componentField" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                </form>
                <DialogFooter>
                  <Button type="submit" form="channelForm">
                    <Loader2 v-if="status === 'pending'" class="mr-2 h-4 w-4 animate-spin" />
                    <span v-else>Create Channel</span>
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
          <div class="space-y-0.5 overflow-y-auto flex-1 pr-2">
            <RouterLink
              v-for="channel in channelStore.channels"
              :key="channel.id"
              :to="`/channel/${channel.id}`"
              class="flex items-center gap-2 px-2 py-1.5 rounded text-sm transition-colors hover:bg-accent/50"
              active-class="bg-accent text-accent-foreground"
            >
              <span class="flex-1 truncate">{{ channel.name }}</span>
            </RouterLink>
          </div>
        </div>
        <div class="mt-6 pt-4 border-t">
          <div class="flex items-center justify-between px-2 py-1 mb-2">
            <h3 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Integrations
            </h3>
          </div>
          <div class="space-y-0.5">
            <RouterLink
              to="/bot-management"
              class="flex items-center gap-2 px-2 py-1.5 rounded text-sm transition-colors hover:bg-accent/50"
              active-class="bg-accent text-accent-foreground"
            >
              <Bot class="h-3 w-3" />
              <span class="flex-1 truncate">Manage Bots</span>
            </RouterLink>
          </div>
        </div>
      </div>
    </SidebarContent>

    <SidebarFooter>
      <NavUser />
    </SidebarFooter>
  </Sidebar>
</template>
