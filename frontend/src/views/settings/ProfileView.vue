<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import SettingsLayout from '@/layouts/SettingsLayout.vue'
import userService from '@/services/user.service'
import { useAuthStore } from '@/stores/auth'
import type { UpdateUserPasswordRequest, UpdateUserRequest } from '@/types/responses/user'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader2 } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import * as z from 'zod'

const authStore = useAuthStore()

const router = useRouter()

const { data, isSuccess } = useQuery({
  queryKey: ['profile'],
  queryFn: () => userService.getCurrentUser(),
})

watch(isSuccess, (newValue) => {
  if (newValue) {
    authStore.setUser(data.value!)
  }
})

const updateProfileSchema = toTypedSchema(
  z.object({
    username: z.string().min(3).optional(),
    email: z.email().optional(),
    first_name: z.string().min(1).optional(),
    last_name: z.string().min(1).optional(),
  }),
)

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

const { handleSubmit, isFieldDirty } = useForm({
  validationSchema: updateProfileSchema,
  initialValues: {
    username: authStore.user?.username,
    email: authStore.user?.email,
    first_name: authStore.user?.first_name,
    last_name: authStore.user?.last_name,
  },
})

const { handleSubmit: handlePasswordSubmit, isFieldDirty: isPasswordFieldDirty } = useForm({
  validationSchema: updatePasswordSchema,
})

const { mutateAsync, status: mutationStatus } = useMutation({
  mutationKey: ['updateProfile'],
  mutationFn: (input: UpdateUserRequest) => userService.updateUserProfile(input),
  onSuccess: (response) => {
    authStore.setUser(response)
    toast.success('Profile updated successfully')
  },
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

const saveProfile = handleSubmit(async (values) => {
  await mutateAsync({
    username: values.username,
    email: values.email,
    first_name: values.first_name,
    last_name: values.last_name,
  })
})

const updatePassword = handlePasswordSubmit(async (values) => {
  await mutatePasswordAsync({
    new_password: values.new_password,
  })
})
</script>

<template>
  <SettingsLayout>
    <div class="flex flex-col h-full">
      <!-- Header -->
      <header
        class="flex items-center justify-between px-6 py-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
      >
        <div>
          <h1 class="text-2xl font-semibold">Profile Settings</h1>
          <p class="text-sm text-muted-foreground">Manage your public profile information</p>
        </div>
      </header>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6">
        <div class="max-w-2xl mx-auto space-y-6">
          <!-- Avatar Section -->
          <Card>
            <CardHeader>
              <CardTitle>Profile Picture</CardTitle>
              <CardDescription>Update your avatar image</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="flex items-center gap-4">
                <Avatar class="h-20 w-20">
                  <AvatarFallback class="text-lg">{{
                    authStore
                      .userDisplayName()
                      .split(' ')
                      .map((n) => n[0])
                      .join('')
                  }}</AvatarFallback>
                </Avatar>
              </div>
            </CardContent>
          </Card>

          <!-- Profile Information -->
          <Card>
            <CardHeader>
              <CardTitle>Profile Information</CardTitle>
              <CardDescription>Update your personal details</CardDescription>
            </CardHeader>
            <CardContent>
              <form class="flex flex-col gap-4" @submit="saveProfile">
                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <FormField
                      v-slot="{ componentField }"
                      name="first_name"
                      :validate-on-blur="!isFieldDirty"
                    >
                      <FormItem>
                        <FormLabel for="first_name">First Name</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="space-y-2">
                    <FormField
                      v-slot="{ componentField }"
                      name="last_name"
                      :validate-on-blur="!isFieldDirty"
                    >
                      <FormItem>
                        <FormLabel for="last_name">Last Name</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>
                </div>
                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <FormField
                      v-slot="{ componentField }"
                      name="username"
                      :validate-on-blur="!isFieldDirty"
                    >
                      <FormItem>
                        <FormLabel for="username">Username</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="space-y-2">
                    <FormField
                      v-slot="{ componentField }"
                      name="email"
                      :validate-on-blur="!isFieldDirty"
                    >
                      <FormItem>
                        <FormLabel for="email">Email</FormLabel>
                        <FormControl>
                          <Input type="email" v-bind="componentField" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>
                </div>

                <Button variant="default" class="self-end">
                  <Loader2 class="animate-spin" v-if="mutationStatus === 'pending'" />
                  <span v-else>Save Changes</span>
                </Button>
              </form>
            </CardContent>
          </Card>

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
                <Button variant="destructive" type="submit" class="self-end">
                  <Loader2 class="animate-spin" v-if="mutationPasswordStatus === 'pending'" />
                  <span v-else>Change Password</span>
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </SettingsLayout>
</template>
