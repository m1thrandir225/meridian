<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { Button } from './ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { FormField, FormItem, FormLabel, FormControl } from './ui/form'
import { Input } from './ui/input'
import * as z from 'zod'
import type { UpdateUserPasswordRequest } from '@/types/responses/user'
import userService from '@/services/user.service'
import { useMutation } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'

const authStore = useAuthStore()
const router = useRouter()

const updatePasswordSchema = toTypedSchema(
  z
    .object({
      new_password: z
        .string()
        .min(8, 'Password must be at least 8 characters long')
        .regex(/[A-Z]/, 'Must contain at least one uppercase letter')
        .regex(/[a-z]/, 'Must contain at least one lowercase letter')
        .regex(/[0-9]/, 'Must contain at least one digit')
        .regex(/[^A-Za-z0-9]/, 'Must contain at least one special character'),
      confirm_password: z.string(),
    })
    .refine((data) => data.new_password === data.confirm_password, {
      message: 'Passwords must match',
      path: ['confirm_password'],
    }),
)

const { handleSubmit: handlePasswordSubmit, isFieldDirty: isPasswordFieldDirty } = useForm({
  validationSchema: updatePasswordSchema,
})

const { mutateAsync: mutatePasswordAsync, status: mutationPasswordStatus } = useMutation({
  mutationKey: ['updatePassword'],
  mutationFn: (input: UpdateUserPasswordRequest) => userService.updateUserPassword(input),
  onSuccess: () => {
    authStore.logout()
    toast.success('Password updated successfully. Please log in again.')
    router.push({ name: 'login' })
  },
})

const updatePassword = handlePasswordSubmit(async (values) => {
  await mutatePasswordAsync({
    new_password: values.new_password,
  })
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Account Security</CardTitle>
      <CardDescription>Change your password</CardDescription>
    </CardHeader>
    <CardContent>
      <form class="flex flex-col gap-4" @submit="updatePassword">
        <FormField
          v-slot="{ componentField }"
          name="new_password"
          :validate-on-blur="!isPasswordFieldDirty"
        >
          <FormItem>
            <FormLabel for="new_password">New Password</FormLabel>
            <FormControl>
              <Input type="new_password" v-bind="componentField" />
            </FormControl>
          </FormItem>
        </FormField>
        <FormField
          v-slot="{ componentField }"
          name="confirm_password"
          :validate-on-blur="!isPasswordFieldDirty"
        >
          <FormItem>
            <FormLabel for="confirm_password">Confirm Password</FormLabel>
            <FormControl>
              <Input type="password" v-bind="componentField" />
            </FormControl>
          </FormItem>
        </FormField>
        <Button variant="default" type="submit" class="self-end">
          <Loader2 class="animate-spin" v-if="mutationPasswordStatus === 'pending'" />
          <span v-else>Change Password</span>
        </Button>
      </form>
    </CardContent>
  </Card>
</template>
