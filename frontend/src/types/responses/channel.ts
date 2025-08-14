import type { Channel } from '@/types/models/channel'

export type CreateChannelRequest = {
  name: string
  topic: string
}

export type CreateChannelResponse = Channel
