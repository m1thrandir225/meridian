<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import authService from '@/services/auth.service'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest } from '@/types/responses/auth'
import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()

const router = useRouter()

const { mutateAsync, status } = useMutation({
  mutationKey: ['login'],
  mutationFn: (input: LoginRequest) => authService.login(input),
  onSuccess: (response) => {
    authStore.login(response)
    router.push({
      name: 'home',
    })
  },
})

const handleSubmit = async (values: { login: string; password: string }) => {
  await mutateAsync({ login: values.login, password: values.password })
}
</script>

<template>
  <AuthLayout>
    <LoginForm :is-loading="status === 'pending'" v-on:submit="handleSubmit" />
  </AuthLayout>
</template>
