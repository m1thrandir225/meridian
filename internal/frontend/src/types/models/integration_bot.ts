export interface IntegrationBot {
  id: string
  service_name: string
  created_at: string
  is_revoked: boolean
  target_channels: string[]
}
