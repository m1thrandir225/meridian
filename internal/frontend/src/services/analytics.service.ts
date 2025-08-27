import config from '@/lib/config'
import { apiRequest } from './api.service'
import {
  type DashboardData,
  type UserGrowthData,
  type MessageVolumeData,
  type ChannelActivityData,
  type TopUserData,
  type ReactionUsageData,
} from '@/types/models/analytics'

const analyticsAPI = `${config.apiUrl}/analytics`

export const analyticsService = {
  getDashboardData: (timeRange: string) =>
    apiRequest<DashboardData>({
      url: `${analyticsAPI}/dashboard`,
      method: 'GET',
      params: { timeRange },
      protected: true,
      headers: undefined,
    }),

  getUserGrowth: (startDate: Date, endDate: Date, interval: string) =>
    apiRequest<UserGrowthData[] | null>({
      url: `${analyticsAPI}/user-growth`,
      method: 'GET',
      params: {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        interval,
      },
      protected: true,
      headers: undefined,
    }),

  getMessageVolume: (startDate: Date, endDate: Date, channelId?: string) =>
    apiRequest<MessageVolumeData[]>({
      url: `${analyticsAPI}/message-volume`,
      method: 'GET',
      params: {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        channelId: channelId ?? '',
      },
      protected: true,
      headers: undefined,
    }),

  getChannelActivity: (startDate: Date, endDate: Date, limit: number) =>
    apiRequest<ChannelActivityData[]>({
      url: `${analyticsAPI}/channel-activity`,
      method: 'GET',
      params: {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        limit: limit.toString(),
      },
      protected: true,
      headers: undefined,
    }),

  getTopUsers: (startDate: Date, endDate: Date, limit: number) =>
    apiRequest<TopUserData[]>({
      url: `${analyticsAPI}/top-users`,
      method: 'GET',
      params: {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        limit: limit.toString(),
      },
      protected: true,
      headers: undefined,
    }),

  getReactionUsage: (startDate: Date, endDate: Date) =>
    apiRequest<ReactionUsageData[]>({
      url: `${analyticsAPI}/reaction-usage`,
      method: 'GET',
      params: {
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
      },
      protected: true,
      headers: undefined,
    }),
}

export default analyticsService
