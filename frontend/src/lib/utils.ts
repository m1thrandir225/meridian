import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getUserDisplayName(firstName: string, lastName: string) {
  return `${firstName} ${lastName}`
}

export function getUserInitials(displayName: string) {
  return displayName
    .split(' ')
    .map((n) => n[0])
    .join('')
}

export function buildInviteURL(inviteCode: string): string {
  const baseURL = window.location.origin
  return `${baseURL}/invites/${inviteCode}`
}

export function extractInviteCode(url: string): string | null {
  const match = url.match(/invites\/([a-zA-Z0-9]+)/)
  return match ? match[1] : null
}
