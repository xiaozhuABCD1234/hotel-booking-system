import api from './client'
import type { ApiResponse, LoginRequest, LoginResponse, RegisterRequest, TokenPair } from '@/types'

export const authApi = {
  login(data: LoginRequest) {
    return api.post<ApiResponse<LoginResponse>>('/auth/login', data)
  },

  register(data: RegisterRequest) {
    return api.post<ApiResponse<LoginResponse>>('/auth/register', data)
  },

  refresh(refreshToken: string) {
    return api.post<ApiResponse<TokenPair>>('/auth/refresh', { refreshToken })
  },

  logout() {
    return api.post<ApiResponse>('/auth/logout')
  },
}
