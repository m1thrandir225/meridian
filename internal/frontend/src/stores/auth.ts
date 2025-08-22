import type { AuthenticationTokens } from '@/types/models/tokens'
import type { User } from '@/types/models/user'
import type { LoginResponse } from '@/types/responses/auth'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore(
  'auth',
  () => {
    const accessToken = ref<string | null>(null)
    const refreshToken = ref<string | null>(null)
    const user = ref<User | null>(null)
    const accessTokenExpirationTime = ref<Date | null>(null)

    const login = (data: LoginResponse) => {
      accessToken.value = data.tokens.access_token
      refreshToken.value = data.tokens.refresh_token
      user.value = data.user
      accessTokenExpirationTime.value = new Date(Date.now() + data.tokens.expires_in * 1000)
    }

    const setUser = (newUser: User) => {
      user.value = newUser
    }

    const setTokens = (tokens: AuthenticationTokens) => {
      accessToken.value = tokens.access_token
      refreshToken.value = tokens.refresh_token
      accessTokenExpirationTime.value = new Date(Date.now() + tokens.expires_in * 1000)
    }

    const logout = () => {
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      accessTokenExpirationTime.value = null
    }

    const userDisplayName = () => {
      if (!user.value) {
        return ''
      }
      return `${user.value?.first_name} ${user.value?.last_name}`
    }

    const checkAuth = () => {
      const now = new Date()

      if (
        !refreshToken.value ||
        !accessTokenExpirationTime.value ||
        now > accessTokenExpirationTime.value
      ) {
        return false
      }

      return true
    }

    return {
      user,
      accessToken,
      refreshToken,
      accessTokenExpirationTime,
      checkAuth,
      userDisplayName,
      logout,
      login,
      setUser,
      setTokens,
    }
  },
  {
    persist: true,
  },
)
