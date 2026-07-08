import api from './client'
import type { ApiResponse, Hotel, HotelSearchParams } from '@/types'

export const hotelApi = {
  /** Public: list hotels with filters */
  list(params?: HotelSearchParams) {
    return api.get<ApiResponse<Hotel[]>>('/hotels', { params })
  },

  /** Public: get hotel by ID */
  getById(id: string) {
    return api.get<ApiResponse<Hotel>>(`/hotels/${id}`)
  },

  /** Admin: create hotel */
  create(data: Partial<Hotel>) {
    return api.post<ApiResponse<Hotel>>('/hotels', data)
  },

  /** Admin: update hotel */
  update(id: string, data: Partial<Hotel>) {
    return api.put<ApiResponse<Hotel>>(`/hotels/${id}`, data)
  },

  /** Admin: delete hotel (soft delete) */
  delete(id: string) {
    return api.delete<ApiResponse>(`/hotels/${id}`)
  },
}
