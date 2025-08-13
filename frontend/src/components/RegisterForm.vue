<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { useForm } from 'vee-validate'
import { FormControl, FormField, FormItem, FormLabel } from './ui/form'

const registerSchema = toTypedSchema(
  z.object({
    username: z.string().min(5),
    password: z
      .string()
      .min(8, 'Password must be at least 8 characters long')
      .regex(/[A-Z]/, 'Must contain at least one uppercase letter')
      .regex(/[a-z]/, 'Must contain at least one lowercase letter')
      .regex(/[0-9]/, 'Must contain at least one digit')
      .regex(/[^A-Za-z0-9]/, 'Must contain at least one special character'),
    email: z.email(),
    first_name: z.string().min(2),
    last_name: z.string().min(2),
  }),
)

const props = defineProps<{
  class?: HTMLAttributes['class']
  isLoading: boolean
  onSubmit: (values: {
    email: string
    password: string
    username: string
    first_name: string
    last_name: string
  }) => Promise<void>
}>()

const { handleSubmit, isFieldDirty } = useForm({
  validationSchema: registerSchema,
})

const formSubmit = handleSubmit((values) => {
  props.onSubmit({
    username: values.username,
    email: values.email,
    first_name: values.first_name,
    last_name: values.last_name,
    password: values.password,
  })
})
</script>

<template>
  <div :class="cn('flex flex-col gap-6', props.class)">
    <Card class="overflow-hidden p-0">
      <CardContent class="grid p-0 md:grid-cols-2">
        <div class="bg-muted relative hidden md:block">
          <img
            src="/register.jpg"
            alt="Image"
            class="absolute inset-0 h-full w-full object-right object-cover dark:brightness-[0.2] dark:grayscale"
          />
        </div>
        <form class="p-6 md:p-8" @submit="formSubmit">
          <div class="flex flex-col gap-6">
            <div class="flex flex-col items-center text-center">
              <h1 class="text-2xl font-bold">Welcome to Meridian</h1>
              <p class="text-muted-foreground text-balance">Create your account</p>
            </div>
            <div class="grid gap-3">
              <FormField
                v-slot="{ componentField }"
                name="first_name"
                :validate-on-blur="!isFieldDirty"
              >
                <FormItem>
                  <FormLabel>First Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Your first name" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <div class="grid gap-3">
              <FormField
                v-slot="{ componentField }"
                name="last_name"
                :validate-on-blur="!isFieldDirty"
              >
                <FormItem>
                  <FormLabel>Last Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Your last name" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <div class="grid gap-3">
              <FormField
                v-slot="{ componentField }"
                name="username"
                :validate-on-blur="!isFieldDirty"
              >
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Your username" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <div class="grid gap-3">
              <FormField v-slot="{ componentField }" name="email" :validate-on-blur="!isFieldDirty">
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" placeholder="Your email" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <div class="grid gap-3">
              <FormField
                v-slot="{ componentField }"
                name="password"
                :validate-on-blur="!isFieldDirty"
              >
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="**********" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <Button type="submit" class="w-full"> Register </Button>

            <div class="text-center text-sm">
              Already have an account?
              <RouterLink to="/login" class="underline underline-offset-4"> Sign in </RouterLink>
            </div>
          </div>
        </form>
      </CardContent>
    </Card>
    <div
      class="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4"
    >
      By clicking continue, you agree to our <a href="#">Terms of Service</a> and
      <a href="#">Privacy Policy</a>.
    </div>
  </div>
</template>
