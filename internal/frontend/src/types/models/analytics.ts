export interface DashboardData {
  totalUsers: number
  newUsers: number
  totalMessages: number
  newMessages: number
  activeChannels: number
  newChannels: number
  totalReactions: number
  newReactions: number
}

export interface UserGrowthData {
  date: string
  newUsers: number
  totalUsers: number
}

export interface MessageVolumeData {
  date: string
  messageCount: number
}

export interface ChannelActivityData {
  channelId: string
  channelName: string
  messageCount: number
  memberCount: number
}

export interface TopUserData {
  userId: string
  username: string
  messageCount: number
}

export interface ReactionUsageData {
  name: string
  value: number
  emoji?: string
}
