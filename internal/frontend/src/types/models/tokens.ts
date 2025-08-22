export interface AuthenticationTokens {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
}

export type RefreshTokenRequest = {
  refresh_token: string
}
