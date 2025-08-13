<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { FormField, FormControl, FormItem, FormLabel } from '@/components/ui/form'
const loginSchema = toTypedSchema(
  z.object({
    login: z.string(),
    password: z.string(),
  }),
)

const props = defineProps<{
  class?: HTMLAttributes['class']
  isLoading: boolean
  onSubmit: (values: { login: string; password: string }) => Promise<void>
}>()

const { handleSubmit, isFieldDirty } = useForm({
  validationSchema: loginSchema,
})

const formSubmit = handleSubmit((values) => {
  props.onSubmit({
    login: values.login,
    password: values.password,
  })
})
</script>

<template>
  <div :class="cn('flex flex-col gap-6', props.class)">
    <Card class="overflow-hidden p-0">
      <CardContent class="grid p-0 md:grid-cols-2">
        <form class="p-6 md:p-8" @submit="formSubmit">
          <div class="flex flex-col gap-6">
            <div class="flex flex-col items-center text-center">
              <h1 class="text-2xl font-bold">Welcome back</h1>
              <p class="text-muted-foreground text-balance">Login to your Meridian account</p>
            </div>
            <div class="grid gap-3">
              <FormField v-slot="{ componentField }" name="login" :validate-on-blur="!isFieldDirty">
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="username or email" v-bind="componentField" />
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
                  <div class="flex items-center">
                    <FormLabel for="password">Password</FormLabel>
                    <RouterLink
                      to="/forgot-password"
                      class="ml-auto text-sm underline-offset-2 hover:underline"
                    >
                      Forgot your password?
                    </RouterLink>
                  </div>
                  <FormControl>
                    <Input type="password" placeholder="**********" v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
            <Button type="submit" class="w-full"> Login </Button>
            <div class="text-center text-sm">
              Don't have an account?
              <RouterLink to="/register" class="underline underline-offset-4"> Sign up </RouterLink>
            </div>
          </div>
        </form>
        <div class="bg-muted relative hidden md:block">
          <img
            src="/login.jpg"
            alt="Image"
            class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
          />
        </div>
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
