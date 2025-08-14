<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { Button } from './ui/button'
import { FormField, FormItem, FormLabel, FormControl } from './ui/form'
import { Input } from './ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import * as z from 'zod'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { useMutation } from '@tanstack/vue-query'
import type { UpdateUserRequest } from '@/types/responses/user'
import userService from '@/services/user.service'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'

const authStore = useAuthStore()

const updateProfileSchema = toTypedSchema(
  z.object({
    username: z.string().min(3).optional(),
    email: z.email().optional(),
    first_name: z.string().min(1).optional(),
    last_name: z.string().min(1).optional(),
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

const { mutateAsync, status: mutationStatus } = useMutation({
  mutationKey: ['updateProfile'],
  mutationFn: (input: UpdateUserRequest) => userService.updateUserProfile(input),
  onSuccess: (response) => {
    authStore.setUser(response)
    toast.success('Profile updated successfully')
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
</script>

<template>
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
            <FormField v-slot="{ componentField }" name="email" :validate-on-blur="!isFieldDirty">
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
</template>
