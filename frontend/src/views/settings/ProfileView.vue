<script setup lang="ts">
import SettingsLayout from '@/layouts/SettingsLayout.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useQuery } from '@tanstack/vue-query'
import userService from '@/services/user.service'

const authStore = useAuthStore()

const { data, error, status, refetch, isSuccess } = useQuery({
  queryKey: ['profile'],
  queryFn: () => userService.getCurrentUser(),
})

watch(isSuccess, (newValue) => {
  if (newValue) {
    authStore.setUser(data.value!)
  }
})

const saveProfile = () => {
  // Save profile logic would go here
  console.log('Saving profile:')
}
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
            <CardContent class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label for="display-name">Display Name</Label>
                  <Input />
                </div>
                <div class="space-y-2">
                  <Label for="username">Username</Label>
                  <Input id="username" placeholder="Your username" />
                </div>
              </div>
              <div class="space-y-2">
                <Label for="email">Email</Label>
                <Input id="email" type="email" placeholder="Your email address" />
              </div>
            </CardContent>
          </Card>

          <!-- Actions -->
          <div class="flex justify-end gap-2">
            <Button variant="outline">Cancel</Button>
            <Button @click="saveProfile">Save Changes</Button>
          </div>
        </div>
      </div>
    </div>
  </SettingsLayout>
</template>
