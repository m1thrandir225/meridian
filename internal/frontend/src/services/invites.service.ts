import config from '@/lib/config'
import type {
  AcceptChannelInviteRequest,
  ChannelCreateInviteRequest,
} from '@/types/responses/channel_invite'
import { apiRequest } from './api.service'
import type { ChannelInvite } from '@/types/models/channel_invite'
import type { Channel } from '@/types/models/channel'

const channelsApiURL = `${config.apiUrl}/messages/channels`
const invitesApiURL = `${config.apiUrl}/messages/invites`

const inviteService = {
  createInvite: (channelId: string, input: ChannelCreateInviteRequest) =>
    apiRequest<ChannelInvite>({
      url: `${channelsApiURL}/${channelId}/invites`,
      method: 'POST',
      protected: true,
      headers: undefined,
      params: undefined,
      data: input,
    }),
  getInvites: (channelId: string) =>
    apiRequest<ChannelInvite[]>({
      url: `${channelsApiURL}/${channelId}/invites`,
      method: 'GET',
      protected: true,
      headers: undefined,
      params: undefined,
    }),
  acceptInvite: (input: AcceptChannelInviteRequest) =>
    apiRequest<Channel>({
      url: `${invitesApiURL}/accept`,
      method: 'POST',
      protected: true,
      headers: undefined,
      params: undefined,
      data: input,
    }),
  deactivateInvite: (inviteId: string) =>
    apiRequest<void>({
      url: `${invitesApiURL}/${inviteId}`,
      method: 'DELETE',
      protected: true,
      headers: undefined,
      params: undefined,
    }),
}

export default inviteService
