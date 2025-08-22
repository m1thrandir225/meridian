<script setup lang="ts">
import RegisterForm from '@/components/RegisterForm.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import authService from '@/services/auth.service'
import type { RegisterRequest } from '@/types/responses/auth'
import { useMutation } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'

const router = useRouter()

const { mutateAsync, status } = useMutation({
  mutationKey: ['register'],
  mutationFn: (input: RegisterRequest) => authService.register(input),
  onSuccess: () => {
    router.push({
      name: 'login',
    })
  },
  onError: (error) => {
    console.log(error)
    toast.error(error.message)
  },
})

const handleSubmit = async (values: {
  username: string
  password: string
  email: string
  first_name: string
  last_name: string
}) => {
  await mutateAsync({
    username: values.username,
    password: values.password,
    first_name: values.first_name,
    last_name: values.last_name,
    email: values.email,
  })
}
</script>

<template>
  <AuthLayout>
    <RegisterForm v-bind:is-loading="status === 'pending'" v-on:submit="handleSubmit" />
  </AuthLayout>
</template>
