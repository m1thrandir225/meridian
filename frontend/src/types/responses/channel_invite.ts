export type ChannelCreateInviteRequest = {
  expires_at: string
  max_uses?: number
}

export type AcceptChannelInviteRequest = {
  invite_code: string
}
