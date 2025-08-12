import type { AuthenticationTokens } from '../models/tokens'
import type { User } from '../models/user'

export type LoginRequest = {
  login: string
  password: string
}

export type LoginResponse = {
  user: User
  tokens: AuthenticationTokens
}

export type RegisterRequest = {
  username: string
  email: string
  first_name: string
  last_name: string
  password: string
}

export type RegisterResponse = User

export type RefreshTokenRequest = {
  refresh_token: string
}

export type RefreshTokenResponse = AuthenticationTokens
