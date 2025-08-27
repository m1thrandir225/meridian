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
import analyticsService from '@/services/analytics.service'
import { useQuery } from '@tanstack/vue-query'
import { AlertCircle, Hash, Heart, MessageSquare, Users } from 'lucide-vue-next'
import { computed, ref } from 'vue'

// Reactive data
const selectedTimeRange = ref('7d')

// Helper function to get start date based on time range
const getStartDate = () => {
  const now = new Date()
  switch (selectedTimeRange.value) {
    case '1h':
      return new Date(now.getTime() - 60 * 60 * 1000)
    case '24h':
      return new Date(now.getTime() - 24 * 60 * 60 * 1000)
    case '7d':
      return new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
    case '30d':
      return new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000)
    default:
      return new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
  }
}

const {
  data: dashboardData,
  status: dashboardDataStatus,
  error: dashboardDataError,
} = useQuery({
  queryKey: ['dashboard-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getDashboardData(selectedTimeRange.value),
  staleTime: 5 * 60 * 1000,
})

const {
  data: userGrowthData,
  status: userGrowthDataStatus,
  error: userGrowthDataError,
} = useQuery({
  queryKey: ['user-growth-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getUserGrowth(getStartDate(), new Date(), 'daily'),
  staleTime: 5 * 60 * 1000,
})

const {
  data: messageVolumeData,
  status: messageVolumeDataStatus,
  error: messageVolumeDataError,
} = useQuery({
  queryKey: ['message-volume-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getMessageVolume(getStartDate(), new Date()),
  staleTime: 5 * 60 * 1000,
})

const {
  data: channelActivityData,
  status: channelActivityDataStatus,
  error: channelActivityDataError,
} = useQuery({
  queryKey: ['channel-activity-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getChannelActivity(getStartDate(), new Date(), 10),
  staleTime: 5 * 60 * 1000,
})

const {
  data: topUsersData,
  status: topUsersDataStatus,
  error: topUsersDataError,
} = useQuery({
  queryKey: ['top-users-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getTopUsers(getStartDate(), new Date(), 10),
  staleTime: 5 * 60 * 1000,
})

const {
  data: reactionUsageData,
  status: reactionUsageDataStatus,
  error: reactionUsageDataError,
} = useQuery({
  queryKey: ['reaction-usage-analytics', selectedTimeRange],
  queryFn: () => analyticsService.getReactionUsage(getStartDate(), new Date()),
  staleTime: 5 * 60 * 1000,
})

const hasErrors = computed(() => {
  return (
    dashboardDataStatus.value === 'error' ||
    userGrowthDataStatus.value === 'error' ||
    messageVolumeDataStatus.value === 'error' ||
    channelActivityDataStatus.value === 'error' ||
    topUsersDataStatus.value === 'error' ||
    reactionUsageDataStatus.value === 'error'
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
      <div v-if="dashboardData" class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Users</CardTitle>
            <Users class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ dashboardData?.total_users || 0 }}</div>
            <p class="text-xs text-muted-foreground">
              +{{ dashboardData?.new_users_today || 0 }} from last period
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Today's messages</CardTitle>
            <MessageSquare class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardData?.messages_today || 0 }}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Active Channels</CardTitle>
            <Hash class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardData?.active_channels || 0 }}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Active Users</CardTitle>
            <Heart class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">
              {{ dashboardData?.active_users || 0 }}
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Loading skeleton for metrics -->
      <div
        v-else-if="dashboardDataStatus === 'pending'"
        class="grid gap-4 md:grid-cols-2 lg:grid-cols-4"
      >
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
              v-if="userGrowthDataStatus === 'pending'"
              class="h-[300px] flex items-center justify-center"
            >
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
            <LineChart
              v-else-if="userGrowthData"
              :data="userGrowthData"
              :categories="['new_users']"
              :index="'period'"
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
              v-if="messageVolumeDataStatus === 'pending'"
              class="h-[300px] flex items-center justify-center"
            >
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
            <BarChart
              v-else-if="messageVolumeData"
              :data="messageVolumeData"
              :categories="['messages']"
              :index="'period'"
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
            <div v-if="channelActivityDataStatus === 'pending'" class="space-y-4">
              <div v-for="i in 5" :key="i" class="flex items-center justify-between">
                <Skeleton class="h-4 w-32" />
                <Skeleton class="h-4 w-20" />
              </div>
            </div>
            <div v-else-if="channelActivityData" class="space-y-4">
              <div
                v-for="channel in channelActivityData"
                :key="channel.channel_id"
                class="flex items-center justify-between"
              >
                <div class="flex items-center space-x-2">
                  <div class="w-2 h-2 rounded-full bg-primary"></div>
                  <span class="text-sm font-medium">{{
                    channel.channel_name || channel.channel_id
                  }}</span>
                </div>
                <div class="flex items-center space-x-2">
                  <span class="text-sm text-muted-foreground"
                    >{{ channel.messages_count }} messages</span
                  >
                  <span class="text-sm text-muted-foreground"
                    >{{ channel.members_count }} members</span
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
            <div v-if="topUsersDataStatus === 'pending'" class="space-y-4">
              <div v-for="i in 5" :key="i" class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                  <Skeleton class="h-6 w-6 rounded-full" />
                  <Skeleton class="h-4 w-20" />
                </div>
                <Skeleton class="h-4 w-16" />
              </div>
            </div>
            <div v-else-if="topUsersData" class="space-y-4">
              <div
                v-for="(user, index) in topUsersData"
                :key="user.user_id"
                class="flex items-center justify-between"
              >
                <div class="flex items-center space-x-2">
                  <div
                    class="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center text-xs font-medium"
                  >
                    {{ index + 1 }}
                  </div>
                  <span class="text-sm font-medium">{{ user.username || user.user_id }}</span>
                </div>
                <span class="text-sm text-muted-foreground">{{ user.messages_sent }} messages</span>
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
              v-if="reactionUsageDataStatus === 'pending'"
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
              v-else-if="reactionUsageData"
              class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4"
            >
              <div
                v-for="reaction in reactionUsageData"
                :key="reaction.reaction_type"
                class="flex flex-col items-center space-y-2 p-4 border rounded-lg"
              >
                <div class="text-2xl">{{ reaction.reaction_type || '👍' }}</div>
                <div class="text-sm font-medium">{{ reaction.percentage }}%</div>
                <div class="text-xs text-muted-foreground">{{ reaction.count }} uses</div>
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
        <Alert v-if="dashboardDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading dashboard data</AlertTitle>
          <AlertDescription>{{ dashboardDataError?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="userGrowthDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading user growth data</AlertTitle>
          <AlertDescription>{{ userGrowthDataError?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="messageVolumeDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading message volume data</AlertTitle>
          <AlertDescription>{{ messageVolumeDataError?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="channelActivityDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading channel activity data</AlertTitle>
          <AlertDescription>{{ channelActivityDataError?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="topUsersDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading top users data</AlertTitle>
          <AlertDescription>{{ topUsersDataError?.message }}</AlertDescription>
        </Alert>

        <Alert v-if="reactionUsageDataStatus === 'error'" variant="destructive">
          <AlertCircle class="h-4 w-4" />
          <AlertTitle>Error loading reaction usage data</AlertTitle>
          <AlertDescription>{{ reactionUsageDataError?.message }}</AlertDescription>
        </Alert>
      </div>
    </div>
  </div>
</template>
