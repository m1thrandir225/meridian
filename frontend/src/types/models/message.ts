export interface Message {
  id: string
  channel_id: string
  sender_user_id?: string
  integration_id?: string
  content_text: string
  parent_message_id?: string
  created_at: string
}

export type MessageCreateRequest = {
  content_text: string
  sender_user_id: string
  is_integration_message?: boolean
  parent_message_id?: string
}
