import { apiRequest } from './api.service'
import config from '@/lib/config'
import type { RefreshTokenRequest } from '@/types/models/tokens'
import {
  type RegisterResponse,
  type LoginRequest,
  type LoginResponse,
  type RegisterRequest,
  type RefreshTokenResponse,
} from '@/types/responses/auth'

const authURL = `${config.apiUrl}/api/v1/auth`

const authService = {
  login: (input: LoginRequest) =>
    apiRequest<LoginResponse>({
      url: `${authURL}/login`,
      method: 'POST',
      params: undefined,
      headers: undefined,
      protected: false,
      data: input,
    }),
  register: (input: RegisterRequest) =>
    apiRequest<RegisterResponse>({
      url: `${authURL}/register`,
      method: 'POST',
      params: undefined,
      headers: undefined,
      protected: false,
      data: input,
    }),
  refreshToken: (input: RefreshTokenRequest) =>
    apiRequest<RefreshTokenResponse>({
      url: `${authURL}/refresh-token`,
      method: 'POST',
      params: undefined,
      headers: undefined,
      protected: false,
      data: input,
    }),
}

export default authService
