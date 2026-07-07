import api from './client'
import type { ApiResponse, Review, CreateReviewRequest, UpdateReviewRequest, PaginatedList, ReviewFull } from '@/types'

export const reviewApi = {
  list(page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<Review>>>('/reviews', { params: { page, pageSize } })
  },

  getById(id: string) {
    return api.get<ApiResponse<Review>>(`/reviews/${id}`)
  },

  create(data: CreateReviewRequest) {
    return api.post<ApiResponse<Review>>('/reviews', data)
  },

  update(id: string, data: UpdateReviewRequest) {
    return api.put<ApiResponse<Review>>(`/reviews/${id}`, data)
  },

  delete(id: string) {
    return api.delete<ApiResponse>(`/reviews/${id}`)
  },

  listByHotel(hotelId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<Review>>>(`/reviews/by-hotel`, { params: { hotelID: hotelId, page, pageSize } })
  },

  listByUser(userId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<Review>>>(`/reviews/by-user`, { params: { userID: userId, page, pageSize } })
  },

  // Report endpoints
  reviewFullByHotel(hotelId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<ReviewFull>>>(`/reports/review-full/by-hotel`, { params: { hotelID: hotelId, page, pageSize } })
  },

  reviewFullByUser(userId: string, page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<ReviewFull>>>(`/reports/review-full/by-user`, { params: { userID: userId, page, pageSize } })
  },
}
