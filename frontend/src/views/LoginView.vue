<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import authService from '@/services/auth.service'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest } from '@/types/responses/auth'
import { useMutation } from '@tanstack/vue-query'

const authStore = useAuthStore()

const { mutateAsync, status } = useMutation({
  mutationKey: ['login'],
  mutationFn: (input: LoginRequest) => authService.login(input),
  onSuccess: (response) => {
    authStore.login(response)
  },
})

const handleSubmit = async () => {
  await mutateAsync({ login: '', password: '' })
}
</script>

<template>
  <AuthLayout>
    <LoginForm :is-loading="status === 'pending'" />
  </AuthLayout>
</template>
