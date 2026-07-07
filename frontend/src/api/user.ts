import api from './client'
import type { ApiResponse, User, CreateUserRequest, UpdateUserRequest, PaginatedList } from '@/types'

export const userApi = {
  list(page = 1, pageSize = 10) {
    return api.get<ApiResponse<PaginatedList<User>>>('/users', { params: { page, pageSize } })
  },

  getById(id: string) {
    return api.get<ApiResponse<User>>(`/users/${id}`)
  },

  create(data: CreateUserRequest) {
    return api.post<ApiResponse<User>>('/users', data)
  },

  update(id: string, data: UpdateUserRequest) {
    return api.put<ApiResponse<User>>(`/users/${id}`, data)
  },

  delete(id: string) {
    return api.delete<ApiResponse>(`/users/${id}`)
  },
}
