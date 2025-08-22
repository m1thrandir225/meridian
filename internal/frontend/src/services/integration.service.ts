import config from '@/lib/config'
import type {
  CreateIntegrationRequest,
  CreateIntegrationResponse,
  ListIntegrationResponse,
  RevokeIntegrationRequest,
  UpdateIntegrationRequest,
  UpdateIntegrationResponse,
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
  revokeIntegration: (input: RevokeIntegrationRequest) =>
    apiRequest<void>({
      method: 'DELETE',
      url: integrationApiURL,
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
