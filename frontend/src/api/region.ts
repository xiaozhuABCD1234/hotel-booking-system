import api from './client'
import type { ApiResponse, Region, GuestBookingStats } from '@/types'

export const regionApi = {
  /** Public: list all regions */
  list() {
    return api.get<ApiResponse<Region[]>>('/regions')
  },

  /** Public: list provinces */
  listProvinces() {
    return api.get<ApiResponse<Region[]>>('/regions/provinces')
  },

  /** Public: list by parent ID */
  listByParent(parentId: number) {
    return api.get<ApiResponse<Region[]>>('/regions/by-parent', { params: { parentID: parentId } })
  },

  /** Public: get by ID */
  getById(id: number) {
    return api.get<ApiResponse<Region>>(`/regions/${id}`)
  },

  /** Admin: create region */
  create(data: Partial<Region>) {
    return api.post<ApiResponse<Region>>('/regions', data)
  },

  /** Admin: update region */
  update(id: number, data: Partial<Region>) {
    return api.put<ApiResponse<Region>>(`/regions/${id}`, data)
  },

  /** Admin: delete region */
  delete(id: number) {
    return api.delete<ApiResponse>(`/regions/${id}`)
  },
}

/** Report API endpoints */
export const reportApi = {
  hotelSummaries(params?: Record<string, unknown>) {
    return api.get<ApiResponse>('/reports/hotel-summaries', { params })
  },

  roomDetails(params?: Record<string, unknown>) {
    return api.get<ApiResponse>('/reports/room-details', { params })
  },

  roomDetailsByHotel(hotelId: string) {
    return api.get<ApiResponse>('/reports/room-details/by-hotel', { params: { hotelID: hotelId } })
  },

  userVipList(params?: Record<string, unknown>) {
    return api.get<ApiResponse>('/reports/user-vip', { params })
  },

  personInfoList() {
    return api.get<ApiResponse>('/reports/person-info')
  },

  guestStats() {
    return api.get<ApiResponse>('/reports/guest-stats')
  },

  topGuests(limit = 10) {
    return api.get<ApiResponse<GuestBookingStats[]>>('/reports/guest-stats/top', { params: { limit } })
  },
}
