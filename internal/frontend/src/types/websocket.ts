export interface WebSocketMessage {
  type: string
  payload: unknown
}

export interface IncomingMessagePayload {
  id: string
  content: string
  sender_user_id?: string
  integration_id?: string
  channel_id: string
  parent_message_id?: string
  timestamp: string
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

export interface IncomingReactionPayload {
  id: string
  message_id: string
  channel_id: string
  user_id: string
  reaction_type: string
  timestamp: string
}

export interface TypingPayload {
  channel_id: string
  user_id: string
}

export interface SendMessagePayload {
  content: string
  channel_id: string
  parent_message_id?: string
}

export interface SendReactionPayload {
  message_id: string
  channel_id: string
  reaction_type: string
}
