/** Login request */
export interface LoginRequest {
  username: string
  password: string
}

/** Register request */
export interface RegisterRequest {
  username: string
  password: string
  phone?: string
  email?: string
  realName?: string
}

/** Auth tokens */
export interface TokenPair {
  accessToken: string
  refreshToken: string
}

/** Login/Register response data (matches backend tokenResponse) */
export interface LoginResponse {
  accessToken: string
  refreshToken: string
}

export interface UserInfo {
  id: string
  username: string
  phone?: string
  email?: string
  role: 'customer' | 'vip' | 'hotel_manager' | 'admin'
  vipLevelId?: number
  points: number
}
