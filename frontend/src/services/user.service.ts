import config from '@/lib/config'
import { apiRequest } from './api.service'
import type { User } from '@/types/models/user'
import {
  type UpdateUserResponse,
  type UpdateUserPasswordRequest,
  type UpdateUserRequest,
} from '@/types/responses/user'

const userServiceAPI = `${config.apiUrl}/auth/me`

const userService = {
  getCurrentUser: () =>
    apiRequest<User>({
      headers: undefined,
      params: undefined,
      method: 'GET',
      protected: true,
      url: userServiceAPI,
      responseType: 'json',
    }),
  updateUserProfile: (input: UpdateUserRequest) =>
    apiRequest<UpdateUserResponse>({
      headers: undefined,
      params: undefined,
      method: 'PUT',
      protected: true,
      url: userServiceAPI,
      responseType: 'json',
      data: input,
    }),
  updateUserPassword: (input: UpdateUserPasswordRequest) =>
    apiRequest({
      method: 'PUT',
      protected: true,
      url: userServiceAPI,
      responseType: 'json',
      data: input,
      headers: undefined,
      params: undefined,
    }),
}

export default userService
