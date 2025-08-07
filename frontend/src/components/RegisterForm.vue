<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
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
        <form class="p-6 md:p-8">
          <div class="flex flex-col gap-6">
            <div class="flex flex-col items-center text-center">
              <h1 class="text-2xl font-bold">Welcome to Meridian</h1>
              <p class="text-muted-foreground text-balance">Create your account</p>
            </div>
            <div class="grid gap-3">
              <Label for="first-name">First Name</Label>
              <Input id="first-name" type="text" placeholder="John" required />
            </div>
            <div class="grid gap-3">
              <Label for="last-name">Last Name</Label>
              <Input id="last-name" type="text" placeholder="Doe" required />
            </div>
            <div class="grid gap-3">
              <Label for="username">Username</Label>
              <Input id="username" type="text" placeholder="johndoe" required />
            </div>
            <div class="grid gap-3">
              <Label for="email">Email</Label>
              <Input id="email" type="email" placeholder="m@example.com" required />
            </div>
            <div class="grid gap-3">
              <Label for="password">Password</Label>
              <Input id="password" type="password" required placeholder="**********" />
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
