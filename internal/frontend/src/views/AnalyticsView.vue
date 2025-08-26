<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { BarChart } from '@/components/ui/chart-bar'
import { LineChart } from '@/components/ui/chart-line'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { useAnalyticsQueries } from '@/queries/analytics'
import { useAuthStore } from '@/stores/auth'
import { AlertCircle, Hash, Heart, MessageSquare, Users } from 'lucide-vue-next'
import { computed, ref } from 'vue'

// Auth store
const authStore = useAuthStore()

// Reactive data
const selectedTimeRange = ref('7d')

// Use TanStack Query for data fetching
const {
  dashboardQuery,
  userGrowthQuery,
  messageVolumeQuery,
  channelActivityQuery,
  topUsersQuery,
  reactionUsageQuery,
} = useAnalyticsQueries(selectedTimeRange.value)

// Computed properties
const hasErrors = computed(() => {
  return (
    dashboardQuery.error ||
    userGrowthQuery.error ||
    messageVolumeQuery.error ||
    channelActivityQuery.error ||
    topUsersQuery.error ||
    reactionUsageQuery.error
  )
})
</script>

<template>
  <div class="min-h-screen bg-background">
    <div class="container mx-auto p-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold tracking-tight">Analytics Dashboard</h1>
          <p class="text-muted-foreground">
            Monitor your platform's performance and user engagement
          </p>
        </div>
        <div class="flex items-center space-x-2">
          <Select v-model="selectedTimeRange" @update:model-value="() => {}">
            <SelectTrigger class="w-[180px]">
              <SelectValue placeholder="Select time range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1h">Last Hour</SelectItem>
              <SelectItem value="24h">Last 24 Hours</SelectItem>
              <SelectItem value="7d">Last 7 Days</SelectItem>
              <SelectItem value="30d">Last 30 Days</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <!-- Key Metrics Cards -->
      <div v-if="dashboardQuery.data" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Users</CardTitle>
            <Users class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ dashboardQuery.data.value?.totalUsers || 0 }}</div>
            <p class="text-xs text-muted-foreground">
              +{{ dashboardQuery.data.value?.newUsers || 0 }} from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Messages Sent</CardTitle>
            <MessageSquare class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardQuery.data.value?.totalMessages || 0 }}
            </div>
            <p class="text-xs text-muted-foreground">
              +{{ dashboardQuery.data.value?.newMessages || 0 }} from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Active Channels</CardTitle>
            <Hash class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardQuery.data.value?.activeChannels || 0 }}
            </div>
            <p class="text-xs text-muted-foreground">
              +{{ dashboardQuery.data.value?.newChannels || 0 }} from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Reactions Given</CardTitle>
            <Heart class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardQuery.data.value?.totalReactions || 0 }}
            </div>
            <p class="text-xs text-muted-foreground">
              +{{ dashboardQuery.data.value?.newReactions || 0 }} from last period
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- Loading skeleton for metrics -->
      <div v-else-if="dashboardQuery.isLoading" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card v-for="i in 4" :key="i">
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <Skeleton class="h-4 w-20" />
            <Skeleton class="h-4 w-4" />
          </CardHeader>
          <CardContent>
            <Skeleton class="h-8 w-16 mb-2" />
            <Skeleton class="h-3 w-24" />
          </CardContent>
        </Card>
      </div>

      <!-- Charts Grid -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <!-- User Growth Chart -->
        <Card class="col-span-4">
          <CardHeader>
            <CardTitle>User Growth</CardTitle>
            <CardDescription>New user registrations over time</CardDescription>
          </CardHeader>
          <CardContent class="pl-2">
            <div
              v-if="userGrowthQuery.isLoading"
              class="h-[300px] flex items-center justify-center"
            >
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
            <LineChart
              v-else-if="userGrowthQuery.data?.value?.data"
              :data="userGrowthQuery.data.value?.data"
              :categories="['newUsers']"
              :index="'date'"
              class="h-[300px]"
            />
            <div v-else class="h-[300px] flex items-center justify-center text-muted-foreground">
              No data available
            </div>
          </CardContent>
        </Card>

        <!-- Message Volume Chart -->
        <Card class="col-span-3">
          <CardHeader>
            <CardTitle>Message Volume</CardTitle>
            <CardDescription>Messages sent per day</CardDescription>
          </CardHeader>
          <CardContent class="pl-2">
            <div
              v-if="messageVolumeQuery.isLoading"
              class="h-[300px] flex items-center justify-center"
            >
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
            <BarChart
              v-else-if="messageVolumeQuery.data?.value?.data"
              :data="messageVolumeQuery.data.value?.data"
              :categories="['messageCount']"
              :index="'date'"
              class="h-[300px]"
            />
            <div v-else class="h-[300px] flex items-center justify-center text-muted-foreground">
              No data available
            </div>
          </CardContent>
        </Card>

        <!-- Channel Activity -->
        <Card class="col-span-4">
          <CardHeader>
            <CardTitle>Most Active Channels</CardTitle>
            <CardDescription>Channels with highest message activity</CardDescription>
          </CardHeader>
          <CardContent>
            <div v-if="channelActivityQuery.isLoading" class="space-y-4">
              <div v-for="i in 5" :key="i" class="flex items-center justify-between">
                <Skeleton class="h-4 w-32" />
                <Skeleton class="h-4 w-20" />
              </div>
            </div>
            <div v-else-if="channelActivityQuery.data?.value?.channels" class="space-y-4">
              <div
                v-for="channel in channelActivityQuery.data.value?.channels"
                :key="channel.channelId"
                class="flex items-center justify-between"
              >
                <div class="flex items-center space-x-2">
                  <div class="w-2 h-2 rounded-full bg-primary"></div>
                  <span class="text-sm font-medium">{{
                    channel.channelName || channel.channelId
                  }}</span>
                </div>
                <div class="flex items-center space-x-2">
                  <span class="text-sm text-muted-foreground"
                    >{{ channel.messageCount }} messages</span
                  >
                  <span class="text-sm text-muted-foreground"
                    >{{ channel.memberCount }} members</span
                  >
                </div>
              </div>
            </div>
            <div v-else class="text-center text-muted-foreground py-8">
              No channel data available
            </div>
          </CardContent>
        </Card>

        <!-- Top Users -->
        <Card class="col-span-3">
          <CardHeader>
            <CardTitle>Top Users</CardTitle>
            <CardDescription>Most active users by messages sent</CardDescription>
          </CardHeader>
          <CardContent>
            <div v-if="topUsersQuery.isLoading" class="space-y-4">
              <div v-for="i in 5" :key="i" class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                  <Skeleton class="h-6 w-6 rounded-full" />
                  <Skeleton class="h-4 w-20" />
                </div>
                <Skeleton class="h-4 w-16" />
              </div>
            </div>
            <div v-else-if="topUsersQuery.data?.value?.users" class="space-y-4">
              <div
                v-for="(user, index) in topUsersQuery.data.value?.users"
                :key="user.userId"
                class="flex items-center justify-between"
              >
                <div class="flex items-center space-x-2">
                  <div
                    class="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center text-xs font-medium"
                  >
                    {{ index + 1 }}
                  </div>
                  <span class="text-sm font-medium">{{ user.username || user.userId }}</span>
                </div>
                <span class="text-sm text-muted-foreground">{{ user.messageCount }} messages</span>
              </div>
            </div>
            <div v-else class="text-center text-muted-foreground py-8">No user data available</div>
          </CardContent>
        </Card>

        <!-- Reaction Usage -->
        <Card class="col-span-7">
          <CardHeader>
            <CardTitle>Reaction Usage</CardTitle>
            <CardDescription>Most popular reactions across the platform</CardDescription>
          </CardHeader>
          <CardContent class="pl-2">
            <div
              v-if="reactionUsageQuery.isLoading"
              class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4"
            >
              <div
                v-for="i in 6"
                :key="i"
                class="flex flex-col items-center space-y-2 p-4 border rounded-lg"
              >
                <Skeleton class="h-8 w-8" />
                <Skeleton class="h-4 w-16" />
                <Skeleton class="h-3 w-12" />
              </div>
            </div>
            <div
              v-else-if="reactionUsageQuery.data?.value?.reactions"
              class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4"
            >
              <div
                v-for="reaction in reactionUsageQuery.data.value?.reactions"
                :key="reaction.name"
                class="flex flex-col items-center space-y-2 p-4 border rounded-lg"
              >
                <div class="text-2xl">{{ reaction.emoji || '👍' }}</div>
                <div class="text-sm font-medium">{{ reaction.name }}</div>
                <div class="text-xs text-muted-foreground">{{ reaction.value }} uses</div>
              </div>
            </div>
            <div v-else class="text-center text-muted-foreground py-8">
              No reaction data available
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Error states -->
      <div v-if="hasErrors" class="space-y-4">
        <Alert v-if="dashboardQuery.error" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading dashboard data</AlertTitle>
          <AlertDescription>{{ dashboardQuery.error.value?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="userGrowthQuery.error" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading user growth data</AlertTitle>
          <AlertDescription>{{ userGrowthQuery.error.value?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="messageVolumeQuery.error" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading message volume data</AlertTitle>
          <AlertDescription>{{ messageVolumeQuery.error.value?.message }}</AlertDescription>
        </Alert>
      </div>
    </div>
  </div>
</template>
