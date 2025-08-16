import config from '@/lib/config'
import type {
  CreateIntegrationRequest,
  CreateIntegrationResponse,
  RevokeIntegrationRequest,
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
}

export default integrationService
