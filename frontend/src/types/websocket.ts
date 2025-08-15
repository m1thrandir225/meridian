export interface WebSocketMessage {
  type: string
  payload: unknown
}

export interface IncomingMessagePayload {
  id: string
  content: string
  sender_id: string
  channel_id: string
  parent_message_id?: string
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
