import config from '@/lib/config'
import type {
  CreateIntegrationRequest,
  CreateIntegrationResponse,
  ListIntegrationResponse,
  RevokeIntegrationRequest,
  UpdateIntegrationRequest,
  UpdateIntegrationResponse,
  UpvokeIntegrationRequest,
  UpvokeIntegrationResponse,
} from '@/types/responses/integration'
import { apiRequest } from './api.service'

const integrationApiURL = `${config.apiUrl}/integrations`

const integrationService = {
  createIntegration: (input: CreateIntegrationRequest) =>
    apiRequest<CreateIntegrationResponse>({
      method: 'POST',
      url: integrationApiURL,
      headers: undefined,
      params: undefined,
      protected: true,
      data: input,
    }),
  deleteIntegration: (integrationId: string) =>
    apiRequest<void>({
      method: 'DELETE',
      url: `${integrationApiURL}/${integrationId}`,
      headers: undefined,
      params: undefined,
      protected: true,
    }),
  revokeIntegration: (input: RevokeIntegrationRequest) =>
    apiRequest<void>({
      method: 'DELETE',
      url: `${integrationApiURL}/revoke`,
      headers: undefined,
      params: undefined,
      protected: true,
      data: input,
    }),
  upvokeIntegration: (input: UpvokeIntegrationRequest) =>
    apiRequest<UpvokeIntegrationResponse>({
      method: 'POST',
      url: `${integrationApiURL}/upvoke`,
      headers: undefined,
      params: undefined,
      protected: true,
      data: input,
    }),
  updateIntegration: (input: UpdateIntegrationRequest) =>
    apiRequest<UpdateIntegrationResponse>({
      method: 'PUT',
      url: `${integrationApiURL}/${input.integration_id}`,
      headers: undefined,
      params: undefined,
      protected: true,
      data: input,
    }),
  listIntegrations: () =>
    apiRequest<ListIntegrationResponse>({
      method: 'GET',
      url: integrationApiURL,
      headers: undefined,
      params: undefined,
      protected: true,
    }),
}

export default integrationService
