<script setup lang="ts">
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import UpdateProfileForm from '@/components/UpdateProfileForm.vue'
import SettingsLayout from '@/layouts/SettingsLayout.vue'
import userService from '@/services/user.service'
import { useAuthStore } from '@/stores/auth'
import { useQuery } from '@tanstack/vue-query'
import { watch } from 'vue'

const authStore = useAuthStore()

const { data, isSuccess } = useQuery({
  queryKey: ['profile'],
  queryFn: () => userService.getCurrentUser(),
})

watch(isSuccess, (newValue) => {
  if (newValue) {
    authStore.setUser(data.value!)
  }
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
              <CardTitle>Profile Details</CardTitle>
              <CardDescription>Your details</CardDescription>
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
                <div class="grid grid-cols-2 gap-2">
                  <div class="grid">
                    <p class="text-sm font-bold">First Name</p>
                    <p class="">{{ authStore.user?.first_name }}</p>
                  </div>
                  <div class="grid">
                    <p class="text-sm font-bold">Last Name</p>
                    <p class="">{{ authStore.user?.last_name }}</p>
                  </div>
                  <div class="grid">
                    <p class="text-sm font-bold">Username</p>
                    <p class="">{{ authStore.user?.username }}</p>
                  </div>
                  <div class="grid">
                    <p class="text-sm font-bold">Email</p>
                    <p class="">{{ authStore.user?.email }}</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
          <UpdateProfileForm />
        </div>
      </div>
    </div>
  </SettingsLayout>
</template>
