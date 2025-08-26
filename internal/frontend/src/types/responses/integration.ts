import type { IntegrationBot } from '../models/integration_bot'

export type CreateIntegrationRequest = {
  service_name: string
  target_channel_ids: string[]
}

export type CreateIntegrationResponse = {
  service_name: string
  target_channels: string[]
  token: string
  token_lookup_hash: string
  is_revoked: boolean
  created_at: string
  id: string
}

export type RevokeIntegrationRequest = {
  integration_id: string
}

export type UpvokeIntegrationRequest = {
  integration_id: string
}

export type UpvokeIntegrationResponse = CreateIntegrationResponse

export type UpdateIntegrationRequest = {
  integration_id: string
  target_channel_ids: string[]
}

export type UpdateIntegrationResponse = IntegrationBot

export type ListIntegrationResponse = {
  integrations: IntegrationBot[]
}
