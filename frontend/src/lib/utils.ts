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
