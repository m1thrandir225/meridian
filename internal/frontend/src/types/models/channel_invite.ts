export interface ChannelInvite {
  id: string
  channel_id: string
  invite_code: string
  expires_at: string
  max_uses?: number
  current_uses: number
  created_at: string
  is_active: boolean
}
