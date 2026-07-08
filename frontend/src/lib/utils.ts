import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

type AxiosErrorShape = {
  response?: {
    data?: {
      message?: string
      error?: {
        message?: string
      }
    }
  }
}

const defaultFallback = '网络错误，请稍后重试'

/**
 * Extract a human-readable error message from a caught error.
 * Prioritizes backend `message` fields for HTTP errors; falls back to
 * `fallback` or a generic network-error string for true connectivity failures.
 */
export function getApiErrorMessage(e: unknown, fallback?: string): string {
  const err = e as AxiosErrorShape

  const msg =
    err?.response?.data?.message ||
    err?.response?.data?.error?.message

  if (msg) return msg

  if (fallback) return fallback

  return defaultFallback
}
