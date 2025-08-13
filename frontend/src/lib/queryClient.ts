import { QueryClient } from '@tanstack/vue-query'
import { isAxiosError } from 'axios'

const MAX_RETRIES = 3
const HTTP_STATUS_TO_NOT_RETRY = [400, 401, 403, 404]

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error) => {
        if (failureCount > MAX_RETRIES) {
          return false
        }
        return !(isAxiosError(error) && HTTP_STATUS_TO_NOT_RETRY.includes(error.response?.status || 0));
      },
    },
  },
})

export default queryClient
