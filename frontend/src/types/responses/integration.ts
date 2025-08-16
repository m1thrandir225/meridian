export type CreateIntegrationRequest = {
  service_name: string
  target_channel_ids: string[]
}

export type CreateIntegrationResponse = {
  service_name: string
  target_channel_ids: string[]
  token: string
  token_lookup_hash: string
  is_revoked: boolean
  created_at: string
  id: string
}

export type RevokeIntegrationRequest = {
  integration_id: string
}

export type WebhookMessageRequest = {
  content_text: string
  target_channel_id?: string
  parent_message_id?: string
  metadata?: Record<string, string>
}

export type CallbackMessageRequest = {
  content_text: string
  target_channel_id?: string
  parent_message_id?: string
  metadata?: Record<string, string>
}
