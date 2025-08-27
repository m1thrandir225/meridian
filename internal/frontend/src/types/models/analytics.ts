export interface DashboardData {
  active_channels: number
  active_users: number
  average_messages_per_user: number
  last_updated: string
  messages_today: number
  new_users_today: number
  peak_hour: number
  total_channels: number
  total_users: number
}

export interface UserGrowthData {
  period: string
  new_users: number
  total_users: number
  growth_rate: number
}

export interface MessageVolumeData {
  avg_length: number
  channels: number
  messages: number
  period: string
}

export interface ChannelActivityData {
  activity_score: number
  channel_id: string
  channel_name: string
  last_message_at: string
  members_count: number
  messages_count: number
}

export interface TopUserData {
  channels_joined: number
  last_active_at: string
  messages_sent: number
  reactions_given: number
  user_id: string
  username: string
}

export interface ReactionUsageData {
  count: number
  percentage: number
  reaction_type: string
}
