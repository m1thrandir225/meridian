<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { buildInviteURL, cn } from '@/lib/utils'
import inviteService from '@/services/invites.service'
import type { ChannelInvite } from '@/types/models/channel_invite'
import type { ChannelCreateInviteRequest } from '@/types/responses/channel_invite'
import { CalendarDate, DateFormatter } from '@internationalized/date'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Copy, Loader2, Trash, Users } from 'lucide-vue-next'
import { toDate } from 'reka-ui/date'
import { useForm } from 'vee-validate'
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'
import * as z from 'zod'
import { Button } from './ui/button'
import { Calendar } from './ui/calendar'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form'
import {
  NumberField,
  NumberFieldContent,
  NumberFieldDecrement,
  NumberFieldIncrement,
  NumberFieldInput,
} from './ui/number-field'
import { PopoverContent, PopoverTrigger } from './ui/popover'
import Popover from './ui/popover/Popover.vue'

interface Props {
  channelId: string
  channelName: string
}

const props = defineProps<Props>()

const createInviteSchema = toTypedSchema(
  z.object({
    maxUse: z.number().optional(),
    expiresAt: z.string(),
  }),
)

const df = new DateFormatter('mk-MK', {
  dateStyle: 'short',
})

const calendarDateToIso = (calendarDate: CalendarDate): string => {
  const date = new Date(calendarDate.year, calendarDate.month - 1, calendarDate.day)
  return date.toISOString()
}

const isoToCalendarDate = (isoString: string): CalendarDate | undefined => {
  try {
    const date = new Date(isoString)
    if (isNaN(date.getTime())) return undefined

    return new CalendarDate(
      date.getFullYear(),
      date.getMonth() + 1, // CalendarDate uses 1-based months
      date.getDate(),
    )
  } catch {
    return undefined
  }
}

const { handleSubmit, setFieldValue, values } = useForm({
  validationSchema: createInviteSchema,
  initialValues: {
    expiresAt: new Date().toISOString(),
    maxUse: undefined,
  },
})

const {
  data: invites,
  status: invitesStatus,
  refetch: refetchInvites,
} = useQuery({
  queryKey: ['invites', props.channelId],
  queryFn: () => inviteService.getInvites(props.channelId),
})

const { mutateAsync: createInvite, status: createInviteStatus } = useMutation({
  mutationKey: ['create-invite', props.channelId],
  mutationFn: (input: ChannelCreateInviteRequest) =>
    inviteService.createInvite(props.channelId, input),
  onSuccess: () => {
    toast.success('Invite created successfully')
    refetchInvites()
  },
  onError: (error) => {
    toast.error('Failed to create invite')
    console.error(error)
  },
})

const { mutateAsync: deactivateInvite, status: deactivateInviteStatus } = useMutation({
  mutationKey: ['deactivate-invite', props.channelId],
  mutationFn: (inviteId: string) => inviteService.deactivateInvite(inviteId),
  onSuccess: () => {
    toast.success('Invite deactivated successfully')
    refetchInvites()
  },
  onError: () => {
    toast.error('Failed to deactivate invite')
  },
})

// State
const showMaxUses = ref(false)
const value = computed({
  get: () => (values.expiresAt ? isoToCalendarDate(values.expiresAt) : undefined),
  set: (val) => val,
})

const placeholder = ref()

const copyInviteLink = async (invite: ChannelInvite) => {
  try {
    const inviteURL = buildInviteURL(invite.invite_code)
    await navigator.clipboard.writeText(inviteURL)
    toast.success('Invite link copied to clipboard')
  } catch (error) {
    console.error('Failed to copy invite link:', error)
    toast.error('Failed to copy invite link')
  }
}

const isExpired = (invite: ChannelInvite) => {
  return new Date(invite.expires_at) < new Date()
}

const handleCreateInvite = handleSubmit(async (values) => {
  await createInvite({
    expires_at: values.expiresAt,
    max_uses: values.maxUse,
  })
})

const handleDeactivateInvite = async (invite: ChannelInvite) => {
  await deactivateInvite(invite.id)
}
</script>

<template>
  <Dialog>
    <DialogTrigger as-child>
      <Button variant="outline" size="sm">
        <Users class="h-4 w-4 mr-2" />
        Invite People
      </Button>
    </DialogTrigger>
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Invite People to {{ channelName }}</DialogTitle>
      </DialogHeader>

      <div class="space-y-4">
        <!-- Create New Invite Form -->
        <form class="space-y-3 p-4 border rounded-lg" @submit.prevent="handleCreateInvite">
          <h3 class="font-medium text-sm">Create New Invite</h3>

          <div class="space-y-2">
            <FormField name="expiresAt">
              <FormItem class="flex flex-col">
                <FormLabel>Expires At</FormLabel>
                <Popover>
                  <PopoverTrigger as-child>
                    <FormControl>
                      <Button
                        variant="outline"
                        :class="
                          cn(
                            'w-[240px] ps-3 text-start font-normal',
                            !value && 'text-muted-foreground',
                          )
                        "
                      >
                        <span>{{ value ? df.format(toDate(value)) : 'Pick a date' }}</span>
                        <CalendarIcon class="ms-auto h-4 w-4 opacity-50" />
                      </Button>
                      <input hidden />
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent class="w-auto p-0">
                    <Calendar
                      v-model:placeholder="placeholder"
                      :model-value="value"
                      calendar-label="Expires At"
                      initial-focus
                      :min-value="new CalendarDate(1900, 1, 1)"
                      :max-value="new CalendarDate(2030, 12, 31)"
                      @update:model-value="
                        (v) => {
                          if (v) {
                            setFieldValue('expiresAt', calendarDateToIso(v as CalendarDate))
                          } else {
                            setFieldValue('expiresAt', undefined)
                          }
                        }
                      "
                    />
                  </PopoverContent>
                </Popover>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>

          <div class="flex items-center space-x-2">
            <Checkbox id="max-uses" v-model="showMaxUses" />
            <Label for="max-uses" class="text-sm">Limit number of uses?</Label>
          </div>

          <div v-if="showMaxUses" class="space-y-2">
            <FormField name="maxUses" v-slot="{ value }">
              <FormItem>
                <FormLabel>Maximum Uses</FormLabel>
                <NumberField
                  class="gap-2"
                  :min="1"
                  :model-value="value"
                  @update:model-value="
                    (v) => {
                      if (v) {
                        setFieldValue('maxUse', v)
                      } else {
                        setFieldValue('maxUse', undefined)
                      }
                    }
                  "
                >
                  <NumberFieldContent>
                    <NumberFieldDecrement />
                    <FormControl>
                      <NumberFieldInput />
                    </FormControl>
                    <NumberFieldIncrement />
                  </NumberFieldContent>
                </NumberField>
              </FormItem>
            </FormField>
          </div>

          <Button type="submit" :disabled="createInviteStatus === 'pending'" class="w-full">
            {{ createInviteStatus === 'pending' ? 'Creating...' : 'Create Invite' }}
          </Button>
        </form>

        <!-- Existing Invites -->
        <div class="space-y-3">
          <h3 class="font-medium text-sm">Active Invites</h3>

          <div
            v-if="invitesStatus === 'success' && invites?.length === 0"
            class="text-sm text-muted-foreground text-center py-4"
          >
            No active invites
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="invite in invites"
              :key="invite.id"
              class="flex items-center justify-between p-3 border rounded-lg"
            >
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <code class="text-xs bg-muted px-2 py-1 rounded">{{ invite.invite_code }}</code>
                  <Badge v-if="isExpired(invite)" variant="destructive" class="text-xs">
                    Expired
                  </Badge>
                  <Badge v-else-if="!invite.is_active" variant="secondary" class="text-xs">
                    Inactive
                  </Badge>
                  <Badge v-else variant="default" class="text-xs"> Active </Badge>
                </div>

                <div class="text-xs text-muted-foreground space-y-1">
                  <div>Expires: {{ df.format(new Date(invite.expires_at)) }}</div>
                  <div v-if="invite.max_uses">
                    Uses: {{ invite.current_uses }} / {{ invite.max_uses }}
                  </div>
                  <div v-else>Uses: {{ invite.current_uses }} (unlimited)</div>
                </div>
              </div>

              <div class="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  @click="copyInviteLink(invite)"
                  :disabled="isExpired(invite) || !invite.is_active"
                >
                  <Copy class="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  @click="handleDeactivateInvite(invite)"
                  :disabled="!invite.is_active || deactivateInviteStatus === 'pending'"
                >
                  <Loader2 v-if="deactivateInviteStatus === 'pending'" class="animate-spin" />
                  <Trash v-else class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
