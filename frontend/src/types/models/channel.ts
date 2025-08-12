export interface Channel {
  id: string
  name: string
  topic: string
  creator_user_id: string
  creation_time: string
  is_archived: boolean
  members_count: number
  last_message_time: string
}

export type ChannelCreateRequest = {
  name: string
  topic: string
  creator_user_id: string
}

export type ChannelJoinRequest = {
  user_id: string
}
