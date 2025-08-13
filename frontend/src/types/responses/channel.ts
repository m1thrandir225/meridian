import type { Channel } from '@/types/models/channel'

export type CreateChannelRequest = {
  name: string
  topic: string
  creator_user_id: string
}

export type CreateChannelResponse = Channel
