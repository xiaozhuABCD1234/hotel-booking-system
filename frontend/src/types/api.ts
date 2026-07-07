/** Unified API response from backend */
export interface ApiResponse<T = unknown> {
  success: boolean
  data?: T
  message?: string
  error?: {
    code: string
    message: string
    details?: unknown
  }
  pagination?: Pagination
  timestamp: string
}

export interface Pagination {
  currentPage: number
  totalPages: number
  totalItems: number
  itemsPerPage: number
  hasNext: boolean
  hasPrev: boolean
}

/** Paginated list response */
export interface PaginatedList<T> {
  items: T[]
  pagination: Pagination
}
