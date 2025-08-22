<script setup lang="ts">
import SettingsLayout from '@/layouts/SettingsLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import UpdatePasswordForm from '@/components/UpdatePasswordForm.vue'
</script>

<template>
  <SettingsLayout>
    <div class="flex flex-col h-full">
      <!-- Header -->
      <header
        class="flex items-center justify-between px-6 py-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
      >
        <div>
          <h1 class="text-2xl font-semibold">Password & Security</h1>
          <p class="text-sm text-muted-foreground">Manage your account security</p>
        </div>
      </header>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6">
        <div class="max-w-2xl mx-auto space-y-6">
          <!-- Change Password -->
          <UpdatePasswordForm />

          <!-- Security Settings -->
          <Card>
            <CardHeader>
              <CardTitle>Security Settings</CardTitle>
              <CardDescription>Additional security options</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-sm font-medium">Account Deletion</h4>
                  <p class="text-sm text-muted-foreground">Permanently delete your account</p>
                </div>
                <Button variant="destructive">Delete Account</Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </SettingsLayout>
</template>
