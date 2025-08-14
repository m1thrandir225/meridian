import config from '@/lib/config'
import type { CreateChannelRequest } from '@/types/responses/channel'
import { apiRequest } from './api.service'
import type { Channel } from '@/types/models/channel'
import type { Message } from '@/types/models/message'

const channelApiURL = `${config.apiUrl}/messages/channels`

const channelService = {
  createChannel: (input: CreateChannelRequest) =>
    apiRequest<Channel>({
      url: channelApiURL,
      method: 'POST',
      headers: undefined,
      params: undefined,
      protected: true,
      data: input,
    }),
  getChannels: () =>
    apiRequest<Channel[]>({
      url: `${channelApiURL}`,
      protected: true,
      headers: undefined,
      params: undefined,
      method: 'GET',
    }),
  getChannel: (input: string) =>
    apiRequest<Channel>({
      url: `${channelApiURL}/${input}`,
      protected: true,
      headers: undefined,
      params: undefined,
      method: 'GET',
    }),
  getMessages: (channelId: string) =>
    apiRequest<Message[]>({
      url: `${channelApiURL}/${channelId}/messages`,
      protected: true,
      headers: undefined,
      params: undefined,
      method: 'GET',
    }),
}

export default channelService
