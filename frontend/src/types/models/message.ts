import type { Reaction } from './reaction'

export interface Message {
  id: string
  channel_id: string
  sender_user_id?: string
  integration_id?: string
  content_text: string
  parent_message_id?: string
  created_at: string
  reactions?: Reaction[]
  sender_user?: {
    id: string
    username: string
    email: string
    first_name: string
    last_name: string
  }
  integration_bot?: {
    id: string
    is_revoked: boolean
    service_name: string
    created_at: string
  }
}
