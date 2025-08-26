import { useQuery } from '@tanstack/vue-query'
import analyticsService from '@/services/analytics.service'

export const useAnalyticsQueries = (timeRange: string) => {
  const getStartDate = () => {
    const now = new Date()
    switch (timeRange) {
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

  const startDate = getStartDate()
  const endDate = new Date()

  // Dashboard overview query
  const dashboardQuery = useQuery({
    queryKey: ['analytics', 'dashboard', timeRange],
    queryFn: () => analyticsService.getDashboardData(timeRange),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  // User growth query
  const userGrowthQuery = useQuery({
    queryKey: ['analytics', 'userGrowth', startDate.toISOString(), endDate.toISOString()],
    queryFn: () => analyticsService.getUserGrowth(startDate, endDate, 'daily'),
    staleTime: 5 * 60 * 1000,
  })

  // Message volume query
  const messageVolumeQuery = useQuery({
    queryKey: ['analytics', 'messageVolume', startDate.toISOString(), endDate.toISOString()],
    queryFn: () => analyticsService.getMessageVolume(startDate, endDate),
    staleTime: 5 * 60 * 1000,
  })

  // Channel activity query
  const channelActivityQuery = useQuery({
    queryKey: ['analytics', 'channelActivity', startDate.toISOString(), endDate.toISOString()],
    queryFn: () => analyticsService.getChannelActivity(startDate, endDate, 10),
    staleTime: 5 * 60 * 1000,
  })

  // Top users query
  const topUsersQuery = useQuery({
    queryKey: ['analytics', 'topUsers', startDate.toISOString(), endDate.toISOString()],
    queryFn: () => analyticsService.getTopUsers(startDate, endDate, 10),
    staleTime: 5 * 60 * 1000,
  })

  // Reaction usage query
  const reactionUsageQuery = useQuery({
    queryKey: ['analytics', 'reactionUsage', startDate.toISOString(), endDate.toISOString()],
    queryFn: () => analyticsService.getReactionUsage(startDate, endDate),
    staleTime: 5 * 60 * 1000,
  })

  return {
    dashboardQuery,
    userGrowthQuery,
    messageVolumeQuery,
    channelActivityQuery,
    topUsersQuery,
    reactionUsageQuery,
    handleTimeRangeChange,
  }
}
