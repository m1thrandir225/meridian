import posthog from 'posthog-js'

export function usePostHog() {
  posthog.init(import.meta.env.VITE_POSTHOG_PUBLIC_KEY, {
    api_host: 'https://eu.i.posthog.com',
    defaults: '2025-05-24',
  })

  return { posthog }
}
